package util

import (
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](i []T) T {
	if len(i) == 0 {
		panic("arg is an empty array/slice")
	}
	var m T
	for idx := 0; idx < len(i); idx++ {
		item := i[idx]
		if idx == 0 {
			m = item
			continue
		}
		if item > m {
			m = item
		}
	}
	return m
}
