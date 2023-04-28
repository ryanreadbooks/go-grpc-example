package service

import (
	"sync"

	"github.com/jinzhu/copier"
)

type InMemoryOrderSaver struct {
	sync.RWMutex
	data map[string]*Orders
}

func NewInMemoryOrderSaver() *InMemoryOrderSaver {
	return &InMemoryOrderSaver{
		data: make(map[string]*Orders),
	}
}

func (s *InMemoryOrderSaver) Save(id string, price float64) error {
	s.Lock()
	defer s.Unlock()

	if o, ok := s.data[id]; ok {
		o.Count += 1
		o.Total += price
	} else {
		s.data[id] = &Orders{Count: 1, Total: price}
	}
	return nil
}

func (s *InMemoryOrderSaver) Get(id string) *Orders {
	s.RLock()
	defer s.RUnlock()

	o, ok := s.data[id]
	if ok {
		var co Orders
		copier.Copy(&co, o)
		return &co
	}
	return nil
}
