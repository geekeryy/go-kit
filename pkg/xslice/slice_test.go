package xslice_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xslice"
)

func TestContains(t *testing.T) {
	log.Println(xslice.Contains([]string{"ab", "c"}, "ab") == true)
	log.Println(xslice.Contains([]string{"ab", "c"}, "a") == false)
	log.Println(xslice.Contains([]string{"ab", "c"}, "c") == true)
	log.Println(xslice.Contains([]int{1, 2, 3}, 2) == true)
	log.Println(xslice.Contains([]int{1, 2, 3}, 4) == false)
	log.Println(xslice.Contains([]uint{1, 2, 3}, uint(2)) == true)
	log.Println(xslice.Contains([]uint{1, 2, 3}, 4) == false)
}

func BenchmarkContainsInt(b *testing.B) {
	var arr []int
	for i := 0; i < 10000; i++ {
		arr = append(arr, i)
	}
	b.Run("Contains any", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.Contains(arr, 1000)
		}
	})

	b.Run("range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range arr {
				if v == 1000 {
					break
				}
			}
		}
	})
}

func BenchmarkContainsString(b *testing.B) {
	var arr []string
	for i := 0; i < 10000; i++ {
		arr = append(arr, strconv.Itoa(i))
	}
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.Contains(arr, "1000")
		}
	})
	b.Run("range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range arr {
				if v == "1000" {
					break
				}
			}
		}
	})
}
