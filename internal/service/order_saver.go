package service

type OrderSaver interface {
	Save(string, float64) error
	Get(string) *Orders
}

type Orders struct {
	Id    string
	Count uint32
	Total float64
}
