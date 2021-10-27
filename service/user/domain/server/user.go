package server

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go-project/common/proto"
	"go-project/service/user/infrastructure/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct{}

func (s *UserServer) List(ctx context.Context, req *proto.UserListRequest) (*proto.UserListResponse, error) {

	dbSpan := opentracing.StartSpan("db", opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))

	users := repository.NewUserRepository().Paginate(int(req.Page), int(req.PageSize))

	dbSpan.Finish()

	rsp := &proto.UserListResponse{}

	for _, user := range users {
		userInfo := proto.UserInfo{
			Id:       user.ID,
			NickName: user.NickName,
		}
		rsp.List = append(rsp.List, &userInfo)
	}

	return rsp, nil
}

func (s *UserServer) Detail(ctx context.Context, req *proto.UserDetailRequest) (*proto.UserDetailResponse, error) {
	user, err := repository.NewUserRepository().GetById(int(req.Id))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	rsp := &proto.UserDetailResponse{
		Id:       user.ID,
		NickName: user.NickName,
	}

	return rsp, nil
}
