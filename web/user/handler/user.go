package handler

import (
	"context"
	"errors"
	_ "errors"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"

	//"github.com/alibaba/sentinel-golang/core/base"
	"github.com/ervin-meng/go-stitch-monster/infrastructure/middleware/logger"
	"github.com/gin-gonic/gin"
	"go-project/common/proto"
	"go-project/web/user/global"
	"go-project/web/user/request"
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
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "參數錯誤"})
			case codes.Internal:
				fallthrough
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务內部錯誤"})
			}
		}

		return
	}
}

func List(ctx *gin.Context) {

	e, b := sentinel.Entry("api-user-list", sentinel.WithTrafficType(base.Inbound))

	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后重试"})
		return
	}

	e.Exit()

	req := request.ListRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e, b = sentinel.Entry("rpc-user-list")

	if b != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务繁忙1，请稍后重试"})
		return
	}

	defer e.Exit()

	tracerCtx, _ := ctx.Get("tracerCtx")

	rsp, err := global.UserServiceClient.List(tracerCtx.(context.Context), &proto.UserListRequest{Page: req.Page, PageSize: req.PageSize})

	if err != nil {
		sentinel.TraceError(e, errors.New("服务繁忙2，请稍后重试"))
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
