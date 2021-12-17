package xslice

import (
	"golang.org/x/exp/constraints"
)

func Contains[T constraints.Ordered](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
