package domain

type Metrics struct {
	ID    string
	MType MType
	Delta *int64
	Value *float64
}
