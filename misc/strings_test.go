package misc

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRamdomString(t *testing.T) {
	fmt.Printf("%s\n", NewRandomString(12).Build())
	fmt.Printf("%s\n", NewRandomString(12).WithNumber(true).Build())
	fmt.Printf("%s\n", NewRandomString(12).WithSpecial(true).Build())
	fmt.Printf("%s\n", NewRandomString(12).WithNumber(false).WithUpper(false).WithLower(true).WithSpecial(true).Build())
	fmt.Printf("%s\n", NewRandomString(12).WithNumber(false).WithSpecial(true).WithLower(false).WithUpper(false).Build())
	fmt.Printf("%s\n", NewRandomString(12).WithNumber(true).WithLower(false).WithUpper(false).WithSpecial(false).Build())
}

func TestRand(t *testing.T) {
	// fmt.Printf("max float64: %f\n", math.MaxFloat64)
	for i := 0; i < 1000; i++ {
		fmt.Printf("%f\n", (20.0+float64(rand.Intn(15)))/10.0)
	}
}
