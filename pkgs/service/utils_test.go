package service

import (
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestMapToSlice(t *testing.T) {
	var (
		count int
		size  int
	)

	size = 10
	m := make(map[int]int)
	s := make([]int, size)
	for count = 0; count < size; count++ {
		m[count] = count + 1
	}

	MapToSlice(m, s)
}
