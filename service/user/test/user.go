package main

import (
	"context"
	"fmt"
	"go-project/common/proto"
)

func TestList() {
	rsp, err := client.List(context.Background(), &proto.UserListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil {
		panic(err)
	}

	for _, user := range rsp.List {
		fmt.Println(user.Id, user.NickName)
	}
}
