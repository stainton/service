package app

import (
	"context"
	"fmt"
	"service-sender/pkg/fakesvc/proto"

	"google.golang.org/grpc"
)

func RunClient(target string) {
	c, err := grpc.NewClient(target)
	if err != nil {
		panic(err)
	}
	fClient := proto.NewFakeServiceClient(c)
	stream, err := fClient.GetStream(context.Background(), nil)
	if err != nil {
		fmt.Printf("err: %+v", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			fmt.Printf("recv err: %+v", err)
			return
		}
		fmt.Printf("resp: %+v\n", resp)
	}
}
