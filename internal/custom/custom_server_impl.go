package custom

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc"
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
