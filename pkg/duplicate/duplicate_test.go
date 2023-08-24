// Package main @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/8/23 09:59
package duplicate

import (
	"log"
	"testing"
	"time"
)

func TestDuplicate(t *testing.T) {
	d := New(time.Second * 5)
	for i := 0; i < 10; i++ {
		go func() {
			d.Add("func1", func() {
				log.Println("func1-1")
			})
			d.Add("func2", func() {
				log.Println("func2-1")
			})
		}()
	}
	time.Sleep(time.Second * 3)
	for i := 0; i < 10; i++ {
		go func() {
			d.Add("func1", func() {
				log.Println("func1-2")
			})
			d.Add("func2", func() {
				log.Println("func2-2")
			})
		}()
	}
	time.Sleep(time.Second * 6)
	for i := 0; i < 10; i++ {
		go func() {
			d.Add("func1", func() {
				log.Println("func1-3")
			})
			d.Add("func2", func() {
				log.Println("func2-3")
			})
		}()
	}
	time.Sleep(time.Second * 6)

}
