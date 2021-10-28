package handler

import (
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/ervin-meng/go-conch/infrastructure/middleware/logger"
	"github.com/gin-gonic/gin"
	"go-project/common/proto"
	"go-project/web/user/global"
	"go-project/web/user/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func HandleServiceErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "內部錯誤"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "參數錯誤"})
			}
		}

		return
	}
}

func List(ctx *gin.Context) {
	e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求过于频繁"})
		return
	}
	e.Exit()
	//发送gRpc请求
	tracerCtx, _ := ctx.Get("tracerCtx")
	rsp, err := global.UserServiceClient.List(tracerCtx.(context.Context), &proto.UserListRequest{Page: 1, PageSize: 10})

	if err != nil {
		logger.Global.Errorw("获取用户列表失败")
		HandleServiceErrorToHttp(err, ctx)
		return
	}

	//将gRpc数据转换成api期望的格式
	result := make([]interface{}, 0)

	for _, value := range rsp.List {

		data := response.ListResponse{
			Id:       value.Id,
			NickName: value.NickName,
		}

		result = append(result, data)
	}

	ctx.JSON(http.StatusOK, result)
}

func Detail(ctx *gin.Context) {
	tracerCtx, _ := ctx.Get("tracerCtx")
	rsp, err := global.UserServiceClient.Detail(tracerCtx.(context.Context), &proto.UserDetailRequest{Id: 1})

	if err != nil {
		logger.Global.Errorw("获取用户详情失败")
		HandleServiceErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.ListResponse{
		Id:       rsp.Id,
		NickName: rsp.NickName,
	})
}
