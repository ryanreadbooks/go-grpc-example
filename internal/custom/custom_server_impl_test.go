package custom_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

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

	go server.Serve(listener)

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
