package main

import (
	"context"
	"encoding/json"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	_ "github.com/alibaba/sentinel-golang/core/base"
	_ "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	_ "github.com/alibaba/sentinel-golang/logging"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/event"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/middleware/logger"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/middleware/tracer"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go-project/common/proto"
	"go-project/web/user/global"
	"go-project/web/user/handler"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var app *gin.Engine

func main() {
	//初始化全局日志
	logger.Init()
	//初始化配置
	InitConfigWithCenter()
	//初始化链路追踪器
	tracer.Init(global.Config.Name)
	//初始化用户服务客户端
	InitUserServiceClient()
	//初始化限流器和熔断器
	InitSentinel()
	//初始化路由器
	InitRouter()
	//初始化请求处理句柄
	InitHandler()
	//初始化注册中心
	//register.Init(register.HTTPService, global.Config.Name, global.Config.IP, global.Config.Port)
	//启动服务
	run()
}

func InitConfigWithCenter() {
	viper.AutomaticEnv()

	configFIleEnv := viper.GetString("API_ENV")

	if configFIleEnv == "" {
		configFIleEnv = "pro"
	}

	configFileName := fmt.Sprintf("config/nacos-%s.yml", configFIleEnv)

	v := viper.New()

	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	nacosConfig := global.NacosConfig{}

	if err := v.Unmarshal(&nacosConfig); err != nil {
		panic(err)
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.IP,
			Port:   uint64(nacosConfig.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "config/tmp/nacos/log",
		CacheDir:            "config/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	content, _ := client.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})

	_ = json.Unmarshal([]byte(content), &global.Config)

	go func() {
		_ = client.ListenConfig(vo.ConfigParam{
			DataId: nacosConfig.DataId,
			Group:  nacosConfig.Group,
			OnChange: func(namespace, group, dataId, data string) {
				_ = json.Unmarshal([]byte(data), &global.Config)
				logger.Global.Info("配置中心文件更新")
				logger.Global.Info(global.Config)
			},
		})
	}()
}

func InitConfig() {
	viper.AutomaticEnv()
	configFIleEnv := viper.GetString("API_ENV")
	if configFIleEnv == "" {
		configFIleEnv = "pro"
	}
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("config/%s-%s.yml", configFilePrefix, configFIleEnv)

	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	global.Config = &global.ApiConfig{}

	if err := v.Unmarshal(global.Config); err != nil {
		panic(err)
	}

	fmt.Println(*global.Config)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.Config)
		fmt.Println(global.Config)
	})
}

func InitUserServiceClientWithLB() {
	consulConfig := global.Config.Consul
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.IP, consulConfig.Port, global.Config.Service.User.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		logger.Global.Fatal("用户服务连接失败")
	}

	global.UserServiceClient = proto.NewUserClient(userConn)
}

func InitUserServiceClientWithCenter() {
	//从注册中心获取用户服务信息
	consulConfig := global.Config.Consul

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulConfig.IP, consulConfig.Port)

	consulClient, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	serviceFilter := fmt.Sprintf("Service == \"%s\"", global.Config.Service.User.Name)

	service, err := consulClient.Agent().ServicesWithFilter(serviceFilter)

	if err != nil {
		panic(err)
	}

	for _, value := range service {
		global.Config.Service.User.IP = value.Address
		global.Config.Service.User.Port = value.Port
		break
	}

	//创建gRpc客户端
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	userConn, e := grpc.Dial(fmt.Sprintf("%s:%d", global.Config.Service.User.IP, global.Config.Service.User.Port), opts...)

	if e != nil {
		panic(e)
	}

	global.UserServiceClient = proto.NewUserClient(userConn)
}

func InitUserServiceClient() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	userConn, e := grpc.Dial(fmt.Sprintf("%s:%d", global.Config.Service.User.IP, global.Config.Service.User.Port), opts...)

	if e != nil {
		panic(e)
	}

	global.UserServiceClient = proto.NewUserClient(userConn)
}

func InitSentinel() {
	_ = sentinel.InitDefault()
	//基于qps限流，Flow controller
	_, _ = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct, //直接计数
			ControlBehavior:        flow.Reject, //匀速通过 //flow.Reject 不匀速
			Threshold:              1,
			StatIntervalInMs:       1000,
		},
	})
}

func InitRouter() {
	app = gin.Default()
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	RouterGroup := app.Group("/v1")
	global.Router = &global.ApiRouter{
		RouterGroup,
	}
}

func InitHandler() {
	userRouter := global.Router.Group("user").Use(func(ctx *gin.Context) {
		rootSpan := opentracing.GlobalTracer().StartSpan(ctx.Request.URL.Path)
		tracerCtx := opentracing.ContextWithSpan(context.Background(), rootSpan)
		ctx.Set("tracerCtx", tracerCtx)
		defer rootSpan.Finish()
		ctx.Next()
	})
	userRouter.GET("list", handler.List)
	userRouter.GET("detail", handler.Detail)
}

func run() {
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", global.Config.Port))
	}()
	//监听信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//触发事件
	event.Trigger(event.ServiceTerm)
}
