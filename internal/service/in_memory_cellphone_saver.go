package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"

	"github.com/ryanreadbooks/go-grpc-example/pb"
)

var (
	ErrAlreadyExist = fmt.Errorf("cellphone uuid exists")
)

// 将cellphone信息保存在内存中
type InMemoryCellphoneSaver struct {
	sync.RWMutex
	storage map[string]*pb.Cellphone
}

func NewInMemoryCellphoneSaver() *InMemoryCellphoneSaver {
	return &InMemoryCellphoneSaver{
		storage: make(map[string]*pb.Cellphone),
	}
}

// *InMemoryCellphoneSaver实现CellphoneSaver接口
func (s *InMemoryCellphoneSaver) Save(ctx context.Context, cellphone *pb.Cellphone) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.storage[cellphone.Id]
	if ok {
		// 已经存在相同id的cellphone信息
		return ErrAlreadyExist
	}

	var copiedCellphone pb.Cellphone
	// deep copy cellphone instance and put into map
	err := copier.Copy(&copiedCellphone, cellphone)
	if err != nil {
		// internal error
		return fmt.Errorf("can not save cellphone into memory: %s", err.Error())
	}

	s.storage[cellphone.Id] = &copiedCellphone

	return nil
}

func (s *InMemoryCellphoneSaver) Size() int32 {
	s.RLock()
	defer s.RUnlock()

	return int32(len(s.storage))
}

func (s *InMemoryCellphoneSaver) Exists(id string) bool {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.storage[id]; ok {
		return true
	}
	return false
}

// 接口实现：查找符合条件的手机信息
func (s *InMemoryCellphoneSaver) Search(condition *pb.FilterCondition) []*pb.Cellphone {
	s.RLock()
	defer s.RUnlock()

	var satisfiedCellphone []*pb.Cellphone

	for _, cellphone := range s.storage {
		if conditionSatisfied(condition, cellphone) {
			var cc pb.Cellphone
			copier.Copy(&cc, cellphone)
			satisfiedCellphone = append(satisfiedCellphone, &cc)
		}
	}
	return satisfiedCellphone
}

// 检查cellphone是否符合条件condition
func conditionSatisfied(condition *pb.FilterCondition, cellphone *pb.Cellphone) bool {
	if condition.MinCpuCore > cellphone.Cpu.Cores {
		return false
	}
	if condition.MinBatteryCapacity > cellphone.Battery.Capacity {
		return false
	}
	if condition.MinRamSize > cellphone.Ram.Value {
		return false
	}
	if condition.MinStorageSize > cellphone.Storage.Value {
		return false
	}

	if len(condition.Brands) != 0 {
		brandMatched := false
		for _, brand := range condition.Brands {
			if brand == cellphone.Brand {
				brandMatched = true
				break
			}
		}
		if !brandMatched {
			return false
		}
	}
	return true
}
