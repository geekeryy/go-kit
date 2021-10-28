package xslice_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xslice"
)

type demo struct {
	name string
	age  int
	arr  []string
}

func TestContains(t *testing.T) {
	log.Println(xslice.Contains([]string{"ab", "c"}, "ab") == true)
	log.Println(xslice.Contains([]string{"ab", "c"}, "a") == false)
	log.Println(xslice.Contains([]string{"ab", "c"}, "c") == true)
	log.Println(xslice.Contains([]int{1, 2, 3}, 2) == true)
	log.Println(xslice.Contains([]int{1, 2, 3}, 4) == false)
	log.Println(xslice.Contains([]uint{1, 2, 3}, uint(2)) == true)
	log.Println(xslice.Contains([]uint{1, 2, 3}, 4) == false)
	log.Println(xslice.Contains([]*demo{{"1", 0, []string{"ok", "ok1"}}, {"2", 0, []string{"ok", "ok1"}}, {"3", 0, nil}}, &demo{"1", 0, []string{"ok", "ok1"}}) == true)
	log.Println(xslice.Contains([]*demo{{"1", 0, []string{"ok", "ok1"}}, {"2", 0, []string{"ok", "ok1"}}, {"3", 0, nil}}, &demo{"1", 0, []string{"ok1", "ok1"}}) == false)
	log.Println(xslice.Contains([]*demo{{"1", 0, []string{"ok", "ok1"}}, {"2", 0, []string{"ok", "ok1"}}, {"3", 0, nil}}, demo{"1", 0, []string{"ok", "ok1"}}) == false)
	log.Println(xslice.Contains([]interface{}{demo{"1", 0, []string{"ok", "ok1"}}, demo{"2", 0, []string{"ok", "ok1"}}, demo{"3", 0, nil}}, demo{"1", 0, []string{"1ok", "ok1"}}) == false)
	log.Println(xslice.Contains([]interface{}{demo{"1", 0, []string{"ok", "ok1"}}, demo{"2", 0, []string{"ok", "ok1"}}, demo{"3", 0, nil}}, demo{"3", 0, nil}) == true)

}

func BenchmarkContainsInt(b *testing.B) {
	var arr []int
	for i := 0; i < 10000; i++ {
		arr = append(arr, i)
	}
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.Contains(arr, 1000)
		}
	})
	b.Run("ContainsInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.ContainsInt(arr, 1000)
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
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.ContainsString(arr, "1000")
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


func BenchmarkContains(b *testing.B) {
	var arr []demo
	for i := 0; i < 100000; i++ {
		arr = append(arr, demo{age: i,name: "xixi"})
	}
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			xslice.Contains(arr, demo{age: 10000,name: "xixi"})
		}
	})
	b.Run("range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range arr {
				if v.age == 10000 && v.name=="xixi" {
					break
				}
			}
		}
	})
}
