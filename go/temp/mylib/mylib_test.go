package mylib_test

import (
	"temp/mylib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonEmptyStrings(t *testing.T) {
	input := []string{"hello", "", "world"}
	expected := []string{"hello", "world"}
	var output = mylib.NonEmptyStrings(input)
	assert.Equal(t, expected, output)

	oriCap := cap(input)
	assert.Equal(t, cap(output), oriCap)
}

func TestResverse(t *testing.T) {
	arr := [...]int{1, 2, 3}
	mylib.Reverse(&arr)
}
