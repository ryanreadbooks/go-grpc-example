package service

import (
	"context"

	"github.com/ryanreadbooks/go-grpc-example/pb"
)

// 接口
type CellphoneSaver interface {
	// 保存一条手机信息
	Save(context.Context, *pb.Cellphone) error
	// 返回已有的手机信息的数量
	Size() int32
	// 检查某个id的手机是否存在
	Exists(string) bool
	// 查找符合条件的手机
	Search(*pb.FilterCondition) []*pb.Cellphone
}
