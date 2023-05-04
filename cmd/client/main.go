package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryanreadbooks/go-grpc-example/internal/sample"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

// 创建客户端
func Dial(target string) *grpc.ClientConn {
	// insecure
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not dial to %s, failure connection: %v\n", target, err)
	}
	return conn
}

func InitCellphoneServiceClient(target string) (pb.CellphoneServiceClient, *grpc.ClientConn) {
	conn := Dial(target)
	return pb.NewCellphoneServiceClient(conn), conn
}

func main() {
	// parse flag options
	target := flag.String("target", "", "the target of the grpc server")
	targetService := flag.String("service", "cellphone", "the target service of the server")

	invokeCreateCellphone := flag.Bool("create-cellphone", true, "invoke CreateCellphone method")
	invokeSearchCellphone := flag.Bool("search-cellphone", false, "invoke SearchCellphone method")
	invokeUploadCellphoneCover := flag.Bool("upload-cellphone-cover", false, "invoke UploadCellphoneCover method")
	uploadedCoverImgFilename := flag.String("cover-filename", "", "uploaded cover image filename")
	invokeBuyCellphone := flag.Bool("buy-cellphone", false, "invoke BuyCellphone method")

	flag.Parse()

	log.Printf("dialing to %s\n", *target)

	if *targetService == "cellphone" {
		var client pb.CellphoneServiceClient
		client, conn := InitCellphoneServiceClient(*target)
		defer conn.Close()
		if *invokeCreateCellphone {
			createCellphone(client)
		}
		if *invokeSearchCellphone {
			searchCellphone(client)
		}
		if *invokeUploadCellphoneCover {
			if *uploadedCoverImgFilename == "" {
				log.Fatal("no cover image filename is specified")
			}
			uploadCellphoneCover(client, *uploadedCoverImgFilename)
		}
		if *invokeBuyCellphone {
			buyCellphone(client)
		}
	} else {
		log.Fatalf("target service '%s' not supported\n", *targetService)
	}
}

// 调用rpc的创建cellphone方法
func createCellphone(client pb.CellphoneServiceClient) (res *pb.CreateCellphoneResponse) {
	req := &pb.CreateCellphoneRequest{
		Cellphone: sample.NewCellphone(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.CreateCellphone(ctx, req)
	if err != nil {
		log.Fatalf("can not create cellphone: %v\n", err)
	}

	log.Printf("cellphone: %s created\n", res.Id)
	return
}

// 调用rpc的查询cellphone方法
func searchCellphone(client pb.CellphoneServiceClient) {
	// 创建查询cellphone的条件
	condition := pb.FilterCondition{
		MinCpuCore:         sample.RandomInt32(1, 6),
		MinRamSize:         sample.RandomInt32(1, 8),
		MinStorageSize:     sample.RandomInt32(100, 1024),
		MinBatteryCapacity: sample.RandomInt32(2500, 8000),
		Brands:             []string{"Apple", "Samsung", "Huawei", "Xiaomi", "OPPO", "VIVO", "Honor", "Pixel"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := client.SearchCellphone(ctx, &condition)
	if err != nil {
		log.Fatalf("can not search cellphone: %v\n", err)
	}
	var cellphones []*pb.Cellphone
	for {
		cellphone, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("can not recv from stream: %v\n", err)
		}
		cellphones = append(cellphones, cellphone)
	}
	if len(cellphones) == 0 {
		log.Println("required cellphone not found")
		return
	}
	for _, c := range cellphones {
		log.Printf("cellphone id: %s\n", c.Id)
	}
}

// 调用rpc的上传图片的方法
func uploadCellphoneCover(client pb.CellphoneServiceClient, imgFile string) {
	// 先创建一条手机信息
	res := createCellphone(client)
	imageType := filepath.Ext(imgFile)
	f, err := os.Open(imgFile)
	if err != nil {
		log.Fatalf("can not open file %s: %v\n", imgFile, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("can not acquire file(%s) stat: %v\n", imgFile, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stream, err := client.UploadCellphoneCover(ctx)
	if err != nil {
		log.Fatalf("can not invoke upload cellphone cover method: %v\n", err)
	}

	// 先发文件信息过去
	metaInfoReq := pb.UploadCellphoneCoverRequest{
		Data: &pb.UploadCellphoneCoverRequest_Meta{
			Meta: &pb.CoverMetaInfo{
				Id:        res.Id,
				ImageType: imageType,
				Size:      uint32(stat.Size()),
			},
		},
	}

	// 发送meta info
	err = stream.Send(&metaInfoReq)
	if err != nil {
		log.Fatalf("can not send file meta info: %v\n", err)
	}

	// 开始不断发送block
	var buf []byte = make([]byte, 4096)

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			// 读完了文件的所有内容
			uploadRes, err := stream.CloseAndRecv()
			if err != nil {
				log.Fatalf("can not close and recv: %v\n", err)
			}
			log.Printf("successfully uploaded %d bytes\n", uploadRes.Size)
			break
		}
		if err != nil {
			log.Fatalf("can not read bytes from file\n")
		}

		blockReq := pb.UploadCellphoneCoverRequest{
			Data: &pb.UploadCellphoneCoverRequest_Block{
				Block: buf[:n],
			},
		}
		err = stream.Send(&blockReq)
		if err != nil {
			log.Fatalf("can not send block request to server: %v\n", err)
		}
	}
}

// 调用rpc的BuyCellphone的方法
func buyCellphone(client pb.CellphoneServiceClient) {
	var createdCellphoneIds []string = make([]string, 0, 5)
	// 先创建一些手机信息
	for i := 0; i < 5; i++ {
		createRes := createCellphone(client)
		createdCellphoneIds = append(createdCellphoneIds, createRes.Id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 调用方法
	stream, err := client.BuyCellphone(ctx)
	if err != nil {
		log.Fatalf("can not invoke BuyCellphone method: %v\n", err)
	}

	// 开另外一个goroutine来接收
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Printf("err when receiving buy cellphone response: %v\n", err)
				runtime.Goexit()
			}
			log.Printf("%s: %.3f\n", response.Id, response.Avg)
		}
	}()

	for i := 0; i < 2; i++ {
		for _, id := range createdCellphoneIds {
			// 发送请求
			req := pb.BuyCellphoneRequest{
				Id:    id,
				Price: sample.RandomFloat64(1000.0, 10000.0),
			}
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("can not send buy request: %v\n", err)
			}
		}
	}
	stream.CloseSend()
	<-waitc
}
