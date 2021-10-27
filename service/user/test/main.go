package main

import (
	"go-project/common/proto"
	"google.golang.org/grpc"
)

var client *rpcClient

func main() {
	client = NewRpcClient()
	defer client.Conn.Close()
	TestList()
}

type rpcClient struct {
	proto.UserClient
	Conn *grpc.ClientConn
}

func NewRpcClient() *rpcClient {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	//opts = append(opts, grpc.WithPerRPCCredentials())

	conn, e := grpc.Dial("127.0.0.1:9501", opts...)

	if e != nil {
		panic(e)
	}

	userClient := proto.NewUserClient(conn)

	return &rpcClient{
		userClient,
		conn,
	}
}
