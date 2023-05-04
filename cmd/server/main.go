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

// 安装一个stream interceptor
type myWrappedStream struct {
	grpc.ServerStream
}

// 实现SendMsg方法
func (s *myWrappedStream) SendMsg(data interface{}) error {
	// TODO 就是在这里实现自定义流拦截器的发送逻辑
	// 这里简单打印日志
	// data是需要往流中发送的数据，我们可以在这里对这个准备发送的数据进行自定义操作
	log.Printf("myWrappedStream send a message(%T): %v\n", data, data)
	return s.ServerStream.SendMsg(data)
}

// 实现RecvMsg方法
func (s *myWrappedStream) RecvMsg(data interface{}) error {
	// TODO 就是在这里实现自定义流拦截器的接收逻辑
	// 这里简单打印日志
	err := s.ServerStream.RecvMsg(data) // 从流中接收数据
	if err != nil {
		return err
	}
	// 调用了RecvMsg从流中接收了数据之后，就可通过data访问到接收的数据内容
	// data的类型是在proto文件中定义的流数据的类型
	// 就可以按照自己的需求进一步处理
	log.Printf("myWrappedStream receive a message(%T): %v\n", data, data)

	return nil
}

func newMyServerStream(s grpc.ServerStream) grpc.ServerStream {
	return &myWrappedStream{s}
}

func installServerStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		log.Printf("actual grpc.ServerStream type: %T\n", ss) // *grpc.serverStream
		return handler(srv, newMyServerStream(ss))
	}
}

func main() {
	cellphoneServiceOn := flag.Bool("cellphone", true, "turn on cellphone service")
	customServiceOn := flag.Bool("custom", false, "turn on custom service")

	flag.Parse()

	listener, err := net.Listen("tcp", "127.0.0.1:9527") // tcp监听
	if err != nil {
		log.Fatal(err)
	}
	// 创建服务器
	// 并且添加流拦截器
	server := grpc.NewServer(
		grpc.StreamInterceptor(installServerStreamInterceptor()),
	)
	defer server.GracefulStop()
	defer listener.Close()

	if *cellphoneServiceOn {
		serverImpl := service.NewCellphoneServiceServer("image/server")
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
