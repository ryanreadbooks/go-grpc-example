package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ryanreadbooks/go-grpc-example/internal/sample"
	"github.com/ryanreadbooks/go-grpc-example/internal/service"
	"github.com/ryanreadbooks/go-grpc-example/pb"
)

func TestInMemoryCellphoneSaverSave(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var name1 = "in-men-save-ok"
	var duplicatedCellphone = sample.NewCellphone()

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
			Name:      "duplicate-cellphone",
			Cellphone: duplicatedCellphone,
			Err:       service.ErrAlreadyExist,
		},
		{
			Name:      name1 + "3",
			Cellphone: sample.NewCellphone(),
			Err:       nil,
		},
	}

	saver := service.NewInMemoryCellphoneSaver()
	saver.Save(ctx, duplicatedCellphone)

	for _, tc := range createCellphoneTestCases {
		t.Run(t.Name(), func(it *testing.T) {
			err := saver.Save(ctx, tc.Cellphone)
			require.ErrorIs(it, err, tc.Err)
		})
	}

}

func TestInMemoryCellphoneSaverSearch(t *testing.T) {
	t.Parallel()

	saver := service.NewInMemoryCellphoneSaver()

	cellphones := []*pb.Cellphone{
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 2.0, Cores: 2},
			Ram:     &pb.RAM{Value: 6, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 1024, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4500},
			Brand:   "Apple",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 2.5, Cores: 1},
			Ram:     &pb.RAM{Value: 1, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 250, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 3000},
			Brand:   "Samsung",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 3.2, Cores: 3},
			Ram:     &pb.RAM{Value: 4, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 250, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 3300},
			Brand:   "Huawei",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 4.8, Cores: 4},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 500, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5467},
			Brand:   "OPPO",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 5.0, Cores: 8},
			Ram:     &pb.RAM{Value: 16, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 128, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4333},
			Brand:   "Xiaomi",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 2.2, Cores: 2},
			Ram:     &pb.RAM{Value: 6, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 256, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 2500},
			Brand:   "VIVO",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 3.2, Cores: 1},
			Ram:     &pb.RAM{Value: 2, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 128, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 1455},
			Brand:   "Honor",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 1.8, Cores: 8},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 512, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5633},
			Brand:   "Pixel",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 3.6, Cores: 16},
			Ram:     &pb.RAM{Value: 8, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 256, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 4600},
			Brand:   "Huawei",
		},
		{
			Id:      uuid.NewString(),
			Cpu:     &pb.CPU{MinGhz: 2.8, Cores: 4},
			Ram:     &pb.RAM{Value: 16, Unit: pb.Unit_UnitGB},
			Storage: &pb.Storage{Value: 512, Unit: pb.Unit_UnitGB},
			Battery: &pb.Battery{Capacity: 5999},
			Brand:   "Xiaomi",
		},
	}

	for _, cellphone := range cellphones {
		err := saver.Save(context.Background(), cellphone)
		require.Nil(t, err)
	}

	// conditions
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

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			phones := saver.Search(&tc.Condition)
			require.Equal(tt, tc.ExepctedNum, len(phones))
		})
	}
}
