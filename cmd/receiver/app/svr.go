package app

import (
	"fmt"
	"net"
	"service-sender/pkg/fakesvc/proto"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedFakeServiceServer
}

func (s *Server) GetStream(req *proto.FakeRequest, stream proto.FakeService_GetStreamServer) error {
	for i := 0; i < 100; i++ {
		resp := &proto.FakeResponse{
			Msg:      "hello world",
			MsgId:    fmt.Sprintf("msgid-%v", i),
			StreamId: "null",
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		fmt.Printf("sent message %v\n", i)
	}
	return nil
}

func NewServer() func() {
	gSvr := grpc.NewServer()
	proto.RegisterFakeServiceServer(gSvr, &Server{})
	return func() {
		lis, err := net.Listen("tcp", ":5678")
		if err != nil {
			panic(err)
		}
		fmt.Println("start listening")
		gSvr.Serve(lis)
	}
}
