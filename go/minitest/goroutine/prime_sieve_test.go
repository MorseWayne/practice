package goroutine

import (
	"fmt"
	"testing"
)

func TestPrimeSeive(t *testing.T) {
	ps := NewPrimeSieve()
	res := ps.Generate(10)
	fmt.Println(res)
	expected := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], v)
		}
	}
}
