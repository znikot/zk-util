package misc

import (
	"fmt"
	"testing"
)

type WeightObj struct {
	Name  string
	Value int
}

func (w WeightObj) Weight() int {
	return w.Value
}
func TestWeight(t *testing.T) {
	slice := []WeightObj{
		{Name: "adventure", Value: 5},
		{Name: "truewords", Value: 6},
	}

	wi := NewWeight(slice)

	for i := 0; i < 100; i++ {
		fmt.Printf("%d\n", wi.NextIndex())
	}
}
