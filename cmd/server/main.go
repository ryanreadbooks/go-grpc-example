package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ryanreadbooks/go-grpc-example/internal/custom"
	"github.com/ryanreadbooks/go-grpc-example/internal/service"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

func main() {
	cellphoneServiceOn := flag.Bool("cellphone", true, "turn on cellphone service")
	customServiceOn := flag.Bool("custom", false, "turn on custom service")

	flag.Parse()

	listener, err := net.Listen("tcp", "127.0.0.1:9527") // tcp监听
	if err != nil {
		log.Fatal(err)
	}
	// 创建服务器
	server := grpc.NewServer()
	defer server.GracefulStop()
	defer listener.Close()

	if *cellphoneServiceOn {
		serverImpl := service.NewCellphoneServiceServer()
		pb.RegisterCellphoneServiceServer(server, serverImpl)
	}
	if *customServiceOn {
		customServerImpl := custom.NewCustomServiceServer()
		pb.RegisterCustomServiceServer(server, customServerImpl)
	}

	log.Printf("server is listening on %s\n", listener.Addr().String())
	log.Printf("cellphone service: %v\n", *cellphoneServiceOn)
	log.Printf("custom service: %v\n", *customServiceOn)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("can not serve: %v\n", err)
	}
}
