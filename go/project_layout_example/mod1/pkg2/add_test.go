package pkg2

import (
	"testing"
	"mod1/pkg1"
)

func TestAdd(t *testing.T) {
	a := pkg1.Add(1, 2)
	if a != 3 {
		t.Errorf("Add(1, 2) = %d; want 3", a)
	}
}
