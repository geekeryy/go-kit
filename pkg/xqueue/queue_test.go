// Package main @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/8/25 17:48
package xqueue

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkQueue(b *testing.B) {
	b.Run("RunParallel", func(b *testing.B) {
		b.ReportAllocs()
		q := NewQueueStack[string]()
		b.SetParallelism(10)
		var str = `123`
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				q.Enqueue(str)
				q.DequeueFIFO()
			}
		})
		fmt.Println(q.Len())
	})

	b.Run("Run", func(b *testing.B) {
		b.ReportAllocs()
		q := NewQueueStack[int]()
		for i := 0; i < b.N; i++ {
			q.Enqueue(1)
		}
	})
}

func TestNotifyOnce(t *testing.T) {
	q := NewQueueStack[int]()
	ch := make(chan struct{})
	runtime.GOMAXPROCS(8)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			q.NotifyOnce(ch, nil)
			wg.Done()
		}()
	}
	wg.Wait()
	stop := time.After(time.Second * 1)
	for {
		select {
		case <-stop:
			return
		case <-ch:
			fmt.Println("ch")
		}
	}
}

func TestQueueConcurrent(t *testing.T) {
	q := NewQueueStack[int]()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < 10000; j++ {
				q.Enqueue(i * j)
				//fmt.Println("Enqueue", j)
			}
			wg.Done()
		}(i)
		//go func(i int) {
		//	for j := 0; j < 100000; {
		//		if _, ok := q.DequeueFIFO(); ok {
		//			j++
		//			//fmt.Println("DequeueFIFO", j)
		//		}
		//	}
		//	wg.Done()
		//}(i)
	}
	wg.Wait()
	if q.Len() != 0 {
		t.Error("error len")
	}
	//runtime.GC()
	fmt.Println("GC")
	for i := 0; i < 30; i++ {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
		fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
		fmt.Printf("\tNumGC = %v\n", m.NumGC)
		time.Sleep(time.Second * 10)
	}

}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func TestQueue(t *testing.T) {
	q := NewQueueStack[int]()

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			q.Enqueue(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println(q.Len())
	v := q.DequeueBatchFILO(10)
	fmt.Println(v)
	fmt.Println(q.DequeueBatchFIFO(10))

	fmt.Println(q.Len())

	fmt.Println(q.DequeueFIFO())
	fmt.Println(q.DequeueFILO())

	fmt.Println(q.DequeueFIFO())
	fmt.Println(q.DequeueFILO())

}

func TestPool(t *testing.T) {
	p := sync.Pool{
		New: func() any {
			return &Node[int]{}
		},
	}
	m := sync.Map{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < 10; i++ {
			m.Store(i, p.Get())
		}
		wg.Done()
	}()
	go func() {
		m.Range(func(key, value any) bool {
			a := &Node[int]{}
			p.Put(a)
			return true
		})
		wg.Done()
	}()
	wg.Wait()

}

func TestQueueNull(t *testing.T) {
	fmt.Println(2 * (1 << 30))
	fmt.Println(2 ^ 1)
}
