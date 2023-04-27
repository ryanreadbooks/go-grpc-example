package service_test

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryanreadbooks/go-grpc-example/internal/sample"
	"github.com/ryanreadbooks/go-grpc-example/internal/service"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

// 测试中运行service server
func runTestCellphoneServiceServer(t *testing.T) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 随机端口监听
	require.Nil(t, err)

	// 创建服务器
	server := grpc.NewServer()

	serverImpl := service.NewCellphoneServiceServer()
	pb.RegisterCellphoneServiceServer(server, serverImpl)

	go server.Serve(listener)

	return server, listener
}

// 创建测试过程中使用的client
func makeTestCellphoneServiceClient(t *testing.T, addr string) pb.CellphoneServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	return pb.NewCellphoneServiceClient(conn)
}

// 测试Create服务
func TestCellphoneServiceImplCreateCellphone(t *testing.T) {
	t.Parallel()

	// 初始化测试的服务端和客户端
	server, listener := runTestCellphoneServiceServer(t)
	go server.Serve(listener)
	defer server.GracefulStop()

	client := makeTestCellphoneServiceClient(t, listener.Addr().String())

	name1 := "client-call-create-cellphone"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	duplicatedCellphone := sample.NewCellphone()
	// 先添加一条记录，以便后序测试重复记录的添加
	request := pb.CreateCellphoneRequest{Cellphone: duplicatedCellphone}
	client.CreateCellphone(ctx, &request)

	invalidCellphone := sample.NewCellphone()
	invalidCellphone.Id = "invalid-uuid"

	emptyIdCellphone := sample.NewCellphone()
	emptyIdCellphone.Id = ""

	var createCellphoneTestCases = []struct {
		Name      string
		Cellphone *pb.Cellphone
		Err       error
	}{
		{
			Name:      name1 + "1",
			Cellphone: sample.NewCellphone(),
			Err:       nil,
		},
		{
			Name:      name1 + "2",
			Cellphone: sample.NewCellphone(),
			Err:       nil,
		},
		{
			Name:      "client-call-create-dup-cellphone",
			Cellphone: duplicatedCellphone,
			Err:       service.ErrAlreadyExist,
		},
		{
			Name:      name1 + "3",
			Cellphone: sample.NewCellphone(),
			Err:       nil,
		},
		{
			Name:      "client-call-create-with-invalid-uuid",
			Cellphone: invalidCellphone,
			Err:       service.ErrUUIDInvalid,
		},
		{
			Name:      "client-call-create-with-emptyid",
			Cellphone: emptyIdCellphone,
			Err:       nil,
		},
	}

	for _, tc := range createCellphoneTestCases {
		t.Run(tc.Name, func(tt *testing.T) {
			req := pb.CreateCellphoneRequest{Cellphone: tc.Cellphone}
			ctxInLoop, cancelInLoop := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancelInLoop()
			res, err := client.CreateCellphone(ctxInLoop, &req)
			if err != nil {
				require.NotNil(tt, tc.Err)
			} else {
				require.Nil(tt, tc.Err)
				require.NotEmpty(tt, res.Id)
			}
		})
	}
}

func TestCellphoneServiceImplWithContext(t *testing.T) {
	// 测试在client侧用context叫停调用过程

	// 初始化测试的服务端和客户端
	server, listener := runTestCellphoneServiceServer(t)
	go server.Serve(listener)
	defer server.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	client := makeTestCellphoneServiceClient(t, listener.Addr().String())

	requset := pb.CreateCellphoneRequest{Cellphone: sample.NewCellphone()}

	_, _ = client.CreateCellphone(ctx, &requset)
}

// 测试查找cellphone
func TestCellphoneServiceImplSearchCellphone(t *testing.T) {
	t.Parallel()

	// 初始化测试的服务端和客户端
	server, listener := runTestCellphoneServiceServer(t)
	go server.Serve(listener)
	defer server.GracefulStop()

	client := makeTestCellphoneServiceClient(t, listener.Addr().String())

	cellphones := []*pb.Cellphone{
		{
			Cpu:     &pb.CPU{MinGhz: 2.0, Cores: 2},
			Ram:     &pb.RAM{Value: 6, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 1024, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4500},
			Brand:   "Apple",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 2.5, Cores: 1},
			Ram:     &pb.RAM{Value: 1, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 250, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 3000},
			Brand:   "Samsung",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 3.2, Cores: 3},
			Ram:     &pb.RAM{Value: 4, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 250, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 3300},
			Brand:   "Huawei",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 4.8, Cores: 4},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 500, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5467},
			Brand:   "OPPO",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 5.0, Cores: 8},
			Ram:     &pb.RAM{Value: 16, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 128, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4333},
			Brand:   "Xiaomi",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 2.2, Cores: 2},
			Ram:     &pb.RAM{Value: 6, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 256, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 2500},
			Brand:   "VIVO",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 3.2, Cores: 1},
			Ram:     &pb.RAM{Value: 2, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 128, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 1455},
			Brand:   "Honor",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 1.8, Cores: 8},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 512, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5633},
			Brand:   "Pixel",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 3.6, Cores: 16},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 256, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4600},
			Brand:   "Huawei",
		},
		{
			Cpu:     &pb.CPU{MinGhz: 2.8, Cores: 4},
			Ram:     &pb.RAM{Value: 16, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 512, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5999},
			Brand:   "Xiaomi",
		},
	}

	// 先创建多条手机信息
	for _, cellphone := range cellphones {
		_, err := client.CreateCellphone(context.Background(), &pb.CreateCellphoneRequest{Cellphone: cellphone})
		require.Nil(t, err)
	}

	// 构建测试用例
	name := "case"
	testCases := []*struct {
		Name        string
		Condition   pb.FilterCondition
		ExepctedNum int
	}{
		{
			Name: name + "1",
			Condition: pb.FilterCondition{
				MinCpuCore:         2,
				MinRamSize:         4,
				MinStorageSize:     100,
				MinBatteryCapacity: 2500,
				Brands:             []string{"Samsung", "Huawei", "Xiaomi"},
			},
			ExepctedNum: 4,
		},
		{
			Name: name + "2",
			Condition: pb.FilterCondition{
				MinCpuCore:         4,
				MinRamSize:         8,
				MinStorageSize:     256,
				MinBatteryCapacity: 4000,
				Brands:             []string{"OPPO", "Huawei", "Xiaomi", "Apple"},
			},
			ExepctedNum: 3,
		},
		{
			Name: name + "3",
			Condition: pb.FilterCondition{
				MinCpuCore:         8,
				MinRamSize:         4,
				MinStorageSize:     512,
				MinBatteryCapacity: 2500,
				Brands:             []string{"Pixel"},
			},
			ExepctedNum: 1,
		},
		{
			Name: name + "4",
			Condition: pb.FilterCondition{
				MinCpuCore:         8,
				MinRamSize:         8,
				MinStorageSize:     128,
				MinBatteryCapacity: 4000,
				Brands:             []string{"Apple", "Samsung", "Huawei", "Xiaomi", "OPPO", "VIVO", "Honor", "Pixel"},
			},
			ExepctedNum: 3,
		},
	}

	// 开始查找符合条件的手机信息
	// client侧接收流式响应
	for _, tc := range testCases {
		stream, err := client.SearchCellphone(context.Background(), &tc.Condition)
		require.Nil(t, err)
		var satisfiedCellphone []*pb.Cellphone
		for {
			resCellphone, err := stream.Recv()
			if err == io.EOF { // 已经接受到末尾了
				break
			}
			fmt.Println(resCellphone.Id)
			require.Nil(t, err)
			satisfiedCellphone = append(satisfiedCellphone, resCellphone)
		}
		require.Equal(t, tc.ExepctedNum, len(satisfiedCellphone))
	}
}

func TestCellphoneServiceImplUploadCellphoneCover(t *testing.T) {
	t.Parallel()

	// 初始化测试的服务端和客户端
	server, listener := runTestCellphoneServiceServer(t)
	go server.Serve(listener)
	defer server.GracefulStop()

	client := makeTestCellphoneServiceClient(t, listener.Addr().String())

	testCases := []*struct {
		Name      string
		ImagePath string
		ErrNil    bool
	}{
		{
			Name:      "success-case-1",
			ImagePath: "../../image/client/apple.jpeg",
			ErrNil:    true,
		},
		{
			Name:      "failure-case-too-large",
			ImagePath: "../../image/client/huawei.jpeg",
			ErrNil:    false,
		},
		{
			Name:      "failure-case-invalid-uuid",
			ImagePath: "../../image/client/apple.jpeg",
			ErrNil:    false,
		},
		{
			Name:      "failure-case-id-not-found",
			ImagePath: "../../image/client/apple.jpeg",
			ErrNil:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			// 创建一条手机信息
			res, err := client.CreateCellphone(context.Background(),
				&pb.CreateCellphoneRequest{Cellphone: sample.NewCellphone()})
			require.Nil(tt, err)

			// 上传手机封面
			coverImageName := tc.ImagePath
			imageType := filepath.Ext(coverImageName)
			f, err := os.Open(coverImageName)
			require.Nil(tt, err)
			defer f.Close()

			stat, err := f.Stat()
			require.Nil(tt, err)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// 得到一个stream类型的返回结果
			// client利用这个stream完成请求的发送和响应的接收
			stream, err := client.UploadCellphoneCover(ctx)
			require.Nil(tt, err)

			// 用stream发送内容
			var cellphoneId string
			if tc.Name == "failure-case-invalid-uuid" {
				cellphoneId = "invalid-uuid"
			} else if tc.Name == "failure-case-id-not-found" {
				cellphoneId = "dce5fc07-d7d1-49fe-aaef-183aa779fce2"
			} else {
				cellphoneId = res.Id
			}

			metaInfoRequest := pb.UploadCellphoneCoverRequest{
				Data: &pb.UploadCellphoneCoverRequest_Meta{
					Meta: &pb.CoverMetaInfo{
						Id:        cellphoneId,
						Size:      uint32(stat.Size()),
						ImageType: imageType,
					},
				},
			}

			// 先把meta info发过去
			err = stream.Send(&metaInfoRequest)
			require.Nil(t, err)

			// 然后开始发图片的字节流内容
			var buf []byte = make([]byte, 4096)

			allErrorNil := true

			for {
				// 不断从文件中读4096 bytes数据
				n, err := f.Read(buf)
				if err == io.EOF {
					uploadResponse, err := stream.CloseAndRecv()
					if err != nil {
						allErrorNil = false
						break
					} else {
						require.Nil(tt, err)
					}
					require.Equal(tt, res.Id, uploadResponse.Id)
					require.EqualValues(tt, stat.Size(), uploadResponse.Size)
					break
				}
				require.Nil(tt, err)

				err = stream.Send(&pb.UploadCellphoneCoverRequest{
					Data: &pb.UploadCellphoneCoverRequest_Block{
						Block: buf[:n],
					},
				})
				if err != nil {
					allErrorNil = false
					break
				}
			}
			require.Equal(tt, tc.ErrNil, allErrorNil)
		})
	}
}
