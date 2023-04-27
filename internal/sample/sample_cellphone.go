package sample

// 负责生成cellphone对象

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ryanreadbooks/go-grpc-example/pb"
)

func NewCellphone() *pb.Cellphone {
	cellphone := &pb.Cellphone{
		Id:              uuid.NewString(), // may panic
		Brand:           randomCellphoneBrand(),
		Cpu:             NewCPU(),
		Ram:             NewRAM(),
		Gpu:             NewGPU(),
		Battery:         NewBattery(),
		Storage:         NewStorage(),
		OperatingSystem: NewOperatingSystem(),
		Screen:          NewScreen(),
		Camera:          NewCamera(),
		CreatedAt:       timestamppb.Now(), // protobuf内部的时间类型
	}
	return cellphone
}

func NewCPU() *pb.CPU {
	return &pb.CPU{
		Manufacturer: randomManufacturer(),
		Cores:        randomPowerOfTwo(8),
		MinGhz:       randomFloat64(1.1, 2.5),
		MaxGhz:       randomFloat64(2.5, 5.5),
	}
}

func NewBattery() *pb.Battery {
	return &pb.Battery{
		Capacity: randomInt32(2000, 6000),
	}
}

func NewRAM() *pb.RAM {
	return &pb.RAM{
		Value:   randomPowerOfTwo(6),
		Unit:    pb.Unit_UnitGB,
		DdrType: randomDDRType(),
	}
}

func NewStorage() *pb.Storage {
	return &pb.Storage{
		Value:       randomInt32(100, 8092),
		Unit:        pb.Unit_UnitGB,
		StorageType: randomStorageType(),
	}
}

func NewOperatingSystem() *pb.OperatingSystem {
	return &pb.OperatingSystem{
		Name:    randomOperatingSystem(),
		Version: randomOperatingSystemVersion(),
	}
}

func NewGPU() *pb.GPU {
	return &pb.GPU{
		Manufacturer: randomManufacturer(),
		Memory:       randomInt32(2, 25),
		MemoryUnit:   pb.Unit_UnitGB,
		MinGhz:       randomFloat64(1.0, 3.0),
		MaxGhz:       randomFloat64(3.0, 6.0),
	}
}

func NewScreen() *pb.Screen {
	return &pb.Screen{
		Size:       randomScreenSize(),
		Resolution: randomScreenResolution(),
	}
}

// NewCamera returns a random camera
func NewCamera() *pb.Camera {
	return &pb.Camera{
		Brand: randomCameraBrand(),
		Spec:  randomCameraSpec(),
	}
}
