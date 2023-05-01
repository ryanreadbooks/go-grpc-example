package custom_test

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ryanreadbooks/go-grpc-example/internal/custom"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

// 测试中运行service server
func runTestCustomServiceServer(t *testing.T) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 随机端口监听
	require.Nil(t, err)

	// 创建服务器
	server := grpc.NewServer()

	serverImpl := custom.NewCustomServiceServer()
	pb.RegisterCustomServiceServer(server, serverImpl)

	return server, listener
}

// 创建测试过程中使用的client
func makeTestCustomServiceClient(t *testing.T, addr string) (pb.CustomServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	return pb.NewCustomServiceClient(conn), conn
}

// 测试metadata的携带场景
func TestMetadataCarrying(t *testing.T) {
	t.Parallel()

	server, listener := runTestCustomServiceServer(t)
	defer server.GracefulStop()
	go func() {
		server.Serve(listener)
	}()

	client, conn := makeTestCustomServiceClient(t, listener.Addr().String())
	defer conn.Close()

	type testCase struct {
		Name string
		Data map[string][]string
	}

	newTestCase := func(name string, key []string, value [][]string) *testCase {
		d := make(map[string][]string)
		require.Equal(t, len(key), len(value))
		for i := 0; i < len(key); i++ {
			d[key[i]] = value[i]
		}
		return &testCase{
			Name: name,
			Data: d,
		}
	}

	testCases := []*testCase{
		newTestCase("case1", []string{"name", "age"}, [][]string{{"ryan"}, {"20"}}),
		newTestCase("case2", []string{"areyouok", "thisisgood"}, [][]string{{"1", "2", "3"}, {"we", "qe", "ee"}}),
		newTestCase("case3", []string{"stream-bin", "well-bin", "yes"}, [][]string{{"1243456"}, {"wer\\12sadf"}, {"no"}}),
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			md := metadata.MD{}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			for k, v := range tc.Data {
				md.Append(k, v...)
				for _, vv := range v {
					ctx = metadata.AppendToOutgoingContext(ctx, k, vv)
				}
			}

			// ctxChild := metadata.NewOutgoingContext(ctx, md) // 也可以使用这个函数来添加metadata

			var resHeader metadata.MD
			var resTrailer metadata.MD
			res, err := client.MetadataCarryTest(ctx,
				&pb.CustomRequest{Id: tc.Name},
				grpc.Header(&resHeader),
				grpc.Trailer(&resTrailer))
			// 注意在Unary RPC中获取trailer需要将trailer作为参数传入
			// grpc.Trailer函数返回的是grpc.CallOption
			// 等RPC返回后，响应的trailer的内容就可以获得
			// 同理：响应的header也可以通过同样的方式获取

			require.Nil(t, err)
			for k, v := range tc.Data {
				require.EqualValues(t, v, res.Metadata[k].Values)
			}
		})
	}
}

// status.Convert函数的使用
func TestGPRCStatusAndCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name    string
		Err     error
		Code    codes.Code
		Message string
	}{
		{"case1", status.Errorf(codes.OK, ""), codes.OK, ""},
		{"case2", status.Errorf(codes.Internal, "internal error"), codes.Internal, "internal error"},
		{"case3", status.Errorf(codes.InvalidArgument, "invalid arg"), codes.InvalidArgument, "invalid arg"},
		{"case4", status.Errorf(codes.Unknown, "unknown"), codes.Unknown, "unknown"},
		{"case5", status.Errorf(codes.AlreadyExists, "already exists"), codes.AlreadyExists, "already exists"},
		{"case6", status.Errorf(codes.NotFound, "not found"), codes.NotFound, "not found"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s := status.Convert(tc.Err)
			require.EqualValues(t, tc.Code, s.Code())
			require.EqualValues(t, tc.Message, s.Message())
		})
	}
}

// 这个函数定义的就是server-side unary拦截器的处理函数
func serverSideUnaryInterceptorHandler(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// 我们可以从info中得到调用的相关信息
	// 完整的调用的rpc的名字（包含包名）
	// 因此，其实是可以根据判断FullMethod来确定是否执行拦截器的逻辑（是否执行拦截）
	fullmethod := info.FullMethod
	if fullmethod != "/pb.CustomService/CallWithUnaryInterceptor" {
		return handler(ctx, req)
	}

	// req是进来的请求，在拦截器这里我们就可以对req的内容进行操作了
	reqV, ok := req.(*pb.SimpleRequest)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "can not convert req to *SimpleRequest")
	}

	reqV.Id += "-hijacked"

	// 我们需要在拦截器中手动调用handler，也就是rpc的目标处理函数，并且得到返回结果res
	// 又或者不对res进行操作，直接返回
	res, err := handler(ctx, reqV)

	// 随后我们就可以根据自己的需求对res进行操作了
	resV, ok := res.(*pb.SimpleResponse)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	resV.Data = fullmethod

	return resV, err

}

// 测试中运行service server with unary interceptor
func runTestCustomServiceServerWithUnaryInterceptor(t *testing.T) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 随机端口监听
	require.Nil(t, err)

	// 创建服务器
	// 并且加上指定的unary 拦截器
	// 这个注册的是全局的unary interceptor
	server := grpc.NewServer(grpc.UnaryInterceptor(serverSideUnaryInterceptorHandler))

	serverImpl := custom.NewCustomServiceServer()
	pb.RegisterCustomServiceServer(server, serverImpl)

	go server.Serve(listener)

	return server, listener
}

// 测试server-side unary interceptor的使用
func TestServerSideUnaryInterceptor(t *testing.T) {
	server, listener := runTestCustomServiceServerWithUnaryInterceptor(t)
	defer server.GracefulStop()

	client, conn := makeTestCustomServiceClient(t, listener.Addr().String())
	defer conn.Close()

	res, err := client.CallWithUnaryInterceptor(context.Background(), &pb.SimpleRequest{
		Id: "hello-unary-interceptor",
	})

	require.Nil(t, err)
	require.EqualValues(t, "/pb.CustomService/CallWithUnaryInterceptor", res.Data)

	res, err = client.CallWithUnaryInterceptor2(context.Background(), &pb.SimpleRequest{
		Id: "hello-unary-interceptor2",
	})
	require.Nil(t, err)
	require.NotEqualValues(t, "/pb.CustomService/CallWithUnaryInterceptor", res.Data)
}

// server-side stream handler
// 参数srv：也就是通过RegisterXXXXServiceServer是传入的那个参数
// 这个拦截器在流的传输过程中只被调用了依次
func serverSideStreamInterceptorHandler(srv interface{},
	stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	// fullmethod := info.FullMethod
	log.Printf("stream interceptor: %T, %v, is client stream: %v, is server stream: %v\n", srv, srv,
		info.IsClientStream,
		info.IsServerStream)

	// 依然要我们自己手动调用RPC服务的处理函数
	err := handler(srv, stream)
	if err != nil {
		log.Printf("rpc failed with err: %v\n", err)
	}
	return err
}

// 测试中运行service server with stream interceptor
func runTestCustomServiceServerWithStreamInterceptor(t *testing.T) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 随机端口监听
	require.Nil(t, err)

	// 创建服务器
	// 拦截器
	// 这个注册的是全局的stream interceptor
	server := grpc.NewServer(grpc.StreamInterceptor(serverSideStreamInterceptorHandler))

	serverImpl := custom.NewCustomServiceServer()
	pb.RegisterCustomServiceServer(server, serverImpl)

	go server.Serve(listener)

	return server, listener
}

func TestServerSideStreamInterceptor(t *testing.T) {
	server, listener := runTestCustomServiceServerWithStreamInterceptor(t)
	go server.Serve(listener)
	defer server.GracefulStop()

	client, conn := makeTestCustomServiceClient(t, listener.Addr().String())
	defer conn.Close()

	stream, err := client.CallWithStreamInterceptor(context.Background())
	require.Nil(t, err)

	waitc := make(chan struct{}) // 需要一个channel来通知退出
	// 单独开一个goroutine来接收数据
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Printf("receiving eof reached: %v\n", err)
				break
			}
			if err != nil {
				log.Printf("receiving err: %v\n", err)
				break
			}
			log.Printf("res.Id=%s, res.Data=%s\n", res.Id, res.Data)
		}
		close(waitc)
	}()

	for i := 0; i < 10; i++ {
		is := strconv.Itoa(i)
		err = stream.Send(&pb.SimpleRequest{
			Id: "id-" + is,
		})
		if err == io.EOF {
			break
		}
	}
	stream.CloseSend()
	<-waitc
}
