package test_test

import (
	"encoding/json"
	"log"
	"math/rand"
	"testing"
)

func init() {
	//rand.Seed(0)
}

func TestGrpc_ReloadConfig(t *testing.T) {
	var count int
	for i := 0; i < 100; i++ {
		n := rand.Intn(10)
		log.Println(n)
		if n == 1 {
			count++
		}
	}
	log.Println("count:", count)
}

func BenchmarkMath(b *testing.B) {
	const n = 60
	const a = 1 << n
	const c = a - 1

	b.Run("取模", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = i * 3333333 % a
		}
	})
	b.Run("与", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = i * 3333333 & c
		}
	})

}

func TestDemo(t *testing.T) {
	a := make([]string, 0)
	b := new([]string)
	var c []string
	log.Println(a == nil, *b == nil, c == nil, len(a), len(*b)) // false true true
	log.Printf("%+v %+v %+v \n", a, *b, c)                      // [] [] []
	log.Printf("%p %p %p \n", &a, b, &c)
}

func m(data interface{}) {
	marshal, err := json.Marshal(data)
	log.Println(string(marshal), err)
}
