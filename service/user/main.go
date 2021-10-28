package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go-project/common/infrastructure/event"
	"go-project/common/infrastructure/middleware/logger"
	"go-project/common/infrastructure/middleware/register"
	"go-project/common/infrastructure/middleware/tracer"
	"go-project/common/proto"
	"go-project/service/user/domain/server"
	"go-project/service/user/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//初始化日志中间件
	logger.Init()
	//初始化配置
	InitConfig()
	//初始化数据库
	InitDb()
	//初始化链路追踪中间件
	tracer.Init(global.Config.Name)
	//获取动态服务端口
	//port := utils.GetPort()
	//获取静态服务端口
	port := global.Config.Port
	//创建服务
	RpcServer := grpc.NewServer(grpc.UnaryInterceptor(tracer.OpenTracingGRPCServerInterceptor()))
	//注册用户服务
	proto.RegisterUserServer(RpcServer, &server.UserServer{})
	//注册到服务发现中心
	InitRegister(RpcServer, port)
	//监听服务
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	go func() {
		err = RpcServer.Serve(lis)

		if err != nil {
			panic("failed to start:" + err.Error())
		}
	}()

	//监听信号量
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	//触发服务退出事件
	event.Trigger(event.ServiceTerm)
}

func InitConfig() {
	viper.AutomaticEnv()
	configFIleEnv := viper.GetString("SERVICE_ENV")
	if configFIleEnv == "" {
		configFIleEnv = "pro"
	}
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("infrastructure/config/%s-%s.yml", configFilePrefix, configFIleEnv)

	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	global.Config = &global.ServiceConfig{}

	if err := v.Unmarshal(global.Config); err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.Config)
		fmt.Println(global.Config)
	})
}

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Config.Mysql.User,
		global.Config.Mysql.Password,
		global.Config.Mysql.Host,
		global.Config.Mysql.Port,
		global.Config.Mysql.Db,
	)
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\n\r", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      gormlogger.Info,
			Colorful:      true,
		},
	)
	global.DB_PROJECT, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}

func InitRegister(RpcServer *grpc.Server, port int) {
	//初始化注册中心
	register.Init(global.Config.Consul.IP, global.Config.Consul.Port)
	//注册健康检查接口
	grpc_health_v1.RegisterHealthServer(RpcServer, health.NewServer())
	//创建服务ID
	id := fmt.Sprintf("%s", uuid.NewV4())
	//注册到中心
	register.ServiceRegister(register.RPCService, id, global.Config.Name, global.Config.IP, port)
	//注册服务注销
	event.RegisterHandler(event.ServiceTerm, func() {
		err := register.ServiceDeregister(id)
		if err != nil {
			logger.Global.Info("RPC服务注销失败：", err)
		} else {
			logger.Global.Info("RPC服务注销成功")
		}
	})
}
