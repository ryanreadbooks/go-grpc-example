package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ryanreadbooks/go-grpc-example/internal/service"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9527") // tcp监听
	if err != nil {
		log.Fatal(err)
	}
	// 创建服务器
	server := grpc.NewServer()
	defer server.GracefulStop()

	serverImpl := service.NewCellphoneServiceServer()
	pb.RegisterCellphoneServiceServer(server, serverImpl)

	log.Printf("server is listening on %s\n", listener.Addr().String())
	server.Serve(listener)
}
