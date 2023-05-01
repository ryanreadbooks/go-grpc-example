package custom

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata" // 用来获取grpc传输的metadata信息

	"github.com/ryanreadbooks/go-grpc-example/internal/service"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

type customServiceServer struct {
	pb.UnimplementedCustomServiceServer
}

func NewCustomServiceServer() pb.CustomServiceServer {
	return &customServiceServer{}
}

func (c *customServiceServer) MetadataCarryTest(ctx context.Context, req *pb.CustomRequest) (*pb.CustomResponse, error) {
	// 可以从context中拿到metadata
	var md metadata.MD
	md, ok := metadata.FromIncomingContext(ctx)

	var res map[string]*pb.StringList = make(map[string]*pb.StringList)

	if ok {
		if err := service.CheckContext(ctx); err != nil {
			return nil, err
		}
		// md里面包含了metadata中的键值对
		// md内的key都是小写的
		for k, v := range md {
			log.Printf("server side: %s: %s\n", k, v)
			if _, ok := res[k]; !ok {
				res[k] = &pb.StringList{}
				res[k].Values = make([]string, 0)
			}
			res[k].Values = append(res[k].Values, v...)
		}
	}

	n := strconv.FormatInt(int64(md.Len()), 10)
	responseMetadata := []string{n}
	// 携带响应的header信息
	grpc.SetHeader(ctx, metadata.MD{"num-metadata-recv-header": responseMetadata})
	// 还可以携带响应的trailer信息
	grpc.SetTrailer(ctx, metadata.MD{"num-metadata-recv-trailer": responseMetadata})

	return &pb.CustomResponse{
		Id:       req.Id,
		Metadata: res,
	}, nil
}

// 测试server-side unary interceptor
func (c *customServiceServer) CallWithUnaryInterceptor(ctx context.Context,
	req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	id := req.Id
	now := time.Now().String()

	return &pb.SimpleResponse{Id: id, Data: now}, nil
}

func (c *customServiceServer) CallWithUnaryInterceptor2(ctx context.Context,
	req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	id := req.Id
	now := time.Now().UTC().String()

	return &pb.SimpleResponse{Id: id, Data: now}, nil
}

// 测试server-side stream interceptor
func (c *customServiceServer) CallWithStreamInterceptor(
	stream pb.CustomService_CallWithStreamInterceptorServer) error {

	var i int64 = 0
	for {
		i++
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Aborted, err.Error())
		}
		err = stream.Send(&pb.SimpleResponse{
			Id:   req.Id,
			Data: "res" + strconv.FormatInt(i, 10),
		})
		if err != nil {
			return status.Errorf(codes.Aborted, fmt.Sprintf("send error: %v", err.Error()))
		}
	}
}
