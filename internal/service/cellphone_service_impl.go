package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

var (
	ErrUUIDInvalid = fmt.Errorf("invalid uuid")
)

const (
	maxCoverImageSizeMB = 1
	MaxCoverImageBytes  = uint32(maxCoverImageSizeMB * 1024 * 1024) // bytes
)

// 实现pb生成的service server接口
type cellphoneServiceServer struct {
	pb.UnimplementedCellphoneServiceServer // 必须嵌入这个由protoc生成的结构体
	saver                                  CellphoneSaver
}

func NewCellphoneServiceServer() pb.CellphoneServiceServer {
	return &cellphoneServiceServer{
		saver: NewInMemoryCellphoneSaver(),
	}
}

// 接口实现：添加一台新手机信息
// Unary RPC
func (c *cellphoneServiceServer) CreateCellphone(ctx context.Context,
	req *pb.CreateCellphoneRequest) (response *pb.CreateCellphoneResponse, err error) {

	defer func() {
		if p := recover(); p != nil {
			// panic occurs
			response = nil
			err = status.Error(codes.Internal, fmt.Sprintf("internal panic: %v", p))
		}
	}()

	cellphone := req.Cellphone

	if cellphone.Id == "" {
		// id为空，赋予一个新的id
		log.Printf("requested uuid is empty, now assigning a new one")
		newId := uuid.NewString()
		cellphone.Id = newId
	}

	// uuid不合法
	if err = CheckUUIDValid(cellphone.Id); err != nil {
		response = nil
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	if err = CheckContext(ctx); err != nil {
		response = nil
		return
	}

	// 保存
	if err = c.saver.Save(ctx, cellphone); err != nil {
		var grpcCode codes.Code
		if errors.Is(err, ErrAlreadyExist) {
			grpcCode = codes.AlreadyExists
		} else {
			grpcCode = codes.Internal
		}
		err = status.Error(grpcCode, err.Error())
		response = nil
		return
	}

	response = &pb.CreateCellphoneResponse{}
	response.Id = cellphone.Id
	err = nil
	log.Printf("cellphone with id: %s saved", response.Id)
	return
}

// 接口实现：查找符合条件的手机
// 参数stream用来返回流式响应
// Server streaming RPC
func (c *cellphoneServiceServer) SearchCellphone(condition *pb.FilterCondition,
	stream pb.CellphoneService_SearchCellphoneServer) error {

	// 找出符合条件的手机
	cellphones := c.saver.Search(condition)
	// 其实已经搜索出符合条件的cellphone后，可以直接用unary rpc就返回过去
	// 这里用server streaming rpc只是为了学习这种方式怎样使用

	for _, cellphone := range cellphones {
		// stream嵌入了grpc.ServerStream结构体，里面有一个context.Context成员
		if err := CheckContext(stream.Context()); err != nil {
			return err
		}
		// 通过调用stream的Send方法来执行流式响应
		if err := stream.Send(cellphone); err != nil {
			return err
		}
	}
	return nil
}

// 接口实现:上传手机封面图片
// 参数stream用来接收请求的数据流,并且负责返回响应
// Client streaming RPC
func (c *cellphoneServiceServer) UploadCellphoneCover(stream pb.CellphoneService_UploadCellphoneCoverServer) error {
	// 客户端第一个数据是一个meta info
	request, err := stream.Recv()
	if err != nil {
		log.Printf("can not receive meta data from stream: %v\n", err)
		return err
	}

	// request在proto中定义成了oneof, 所以有两个内容, 类似与union
	cellphoneId := request.GetMeta().Id
	imgSize := request.GetMeta().Size
	imgType := request.GetMeta().ImageType

	// uuid不合法
	if err := CheckUUIDValid(cellphoneId); err != nil {
		log.Printf("cellphone with invalid uuid: %s\n", cellphoneId)
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// 指定的cellphone id不存在
	if !c.saver.Exists(cellphoneId) {
		log.Printf("cellphone with id: %s not found\n", cellphoneId)
		return status.Errorf(codes.NotFound, fmt.Sprintf("cellphone with %s not found", cellphoneId))
	}

	// 文件大小太大
	if imgSize > MaxCoverImageBytes {
		log.Printf("server side imgSize of %d is too large\n", imgSize)
		return status.Errorf(codes.OutOfRange, fmt.Sprintf("provided cover image is larger than %d MB", maxCoverImageSizeMB))
	}
	imgFileName := fmt.Sprintf("../../image/server/%s%s", cellphoneId, imgType)
	imgFile, err := os.Create(imgFileName)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	defer imgFile.Close()

	totalSize := 0

	// 剩下的在for循环中不断接收,接收图片的数据流
	for {
		if err := CheckContext(stream.Context()); err != nil {
			return err
		}
		request, err := stream.Recv()
		// 请求数据接收完成
		if err == io.EOF {
			// 返回响应
			return stream.SendAndClose(&pb.UploadCellphoneCoverResponse{
				Id:   cellphoneId,
				Size: uint32(totalSize),
			})
		}

		if err != nil {
			return err
		}
		// 将接收到的内容写到文件里面
		block := request.GetBlock()
		n, err := imgFile.Write(block)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		totalSize += n
		log.Printf("written %d bytes into %s\n", n, imgFileName)
	}
}