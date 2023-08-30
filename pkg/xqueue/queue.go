// Package queue @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/8/25 17:47
package xqueue

import (
	"sync/atomic"
	"unsafe"
)

type Node[T any] struct {
	value T
	next  *Node[T]
}

type QueueStack[T any] struct {
	head   *Node[T]
	tail   *Node[T]
	count  int32
	notify int32
}

func NewQueueStack[T any]() *QueueStack[T] {
	node := &Node[T]{}
	return &QueueStack[T]{
		head:  node,
		tail:  node,
		count: 0,
	}
}

func (qs *QueueStack[T]) Enqueue(value T) {
	newNode := &Node[T]{value: value}
	for {
		tail := qs.loadTail()
		next := qs.loadNext(tail)
		if tail == qs.loadTail() {
			if next == nil {
				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&tail.next)),
					unsafe.Pointer(next),
					unsafe.Pointer(newNode)) {
					atomic.AddInt32(&qs.count, 1)
					atomic.CompareAndSwapPointer(
						(*unsafe.Pointer)(unsafe.Pointer(&qs.tail)),
						unsafe.Pointer(tail),
						unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&qs.tail)),
					unsafe.Pointer(tail),
					unsafe.Pointer(next))
			}
		}
	}
}

// DequeueFIFO FIFO Dequeue
func (qs *QueueStack[T]) DequeueFIFO() (T, bool) {
	var zero T
	for {
		head := qs.loadHead()
		tail := qs.loadTail()
		next := qs.loadNext(head)
		if head == qs.loadHead() {
			if head == tail {
				if next == nil {
					return zero, false
				}
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&qs.tail)),
					unsafe.Pointer(tail),
					unsafe.Pointer(next))
			} else {
				if next == nil {
					return zero, false
				}
				val := next.value
				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&qs.head)),
					unsafe.Pointer(head),
					unsafe.Pointer(next)) {
					atomic.AddInt32(&qs.count, -1)
					return val, true
				}
			}
		}
	}
}

// DequeueFILO FILO Dequeue
func (qs *QueueStack[T]) DequeueFILO() (T, bool) {
	var prev *Node[T]
	var last *Node[T]
	var zero T
	for {
		head := qs.loadHead()
		last = head
		if last == qs.loadTail() {
			if last == head {
				return zero, false
			}
		}
		for last.next != nil {
			prev = last
			last = last.next
		}
		if prev == nil {
			return zero, false
		}
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&prev.next)),
			unsafe.Pointer(last),
			nil) {
			atomic.AddInt32(&qs.count, -1)
			return last.value, true
		}
	}
}

func (qs *QueueStack[T]) DequeueBatchFIFO(n int) []T {
	var batch []T
	for i := 0; i < n; i++ {
		val, ok := qs.DequeueFIFO()
		if !ok {
			break
		}
		batch = append(batch, val)
	}
	return batch
}

func (qs *QueueStack[T]) DequeueBatchFILO(n int) []T {
	var batch []T

	for i := 0; i < n; i++ {
		val, ok := qs.DequeueFILO()
		if !ok {
			break
		}

		batch = append(batch, val)
	}
	return batch
}

func (qs *QueueStack[T]) NotifyOnce(ch chan<- struct{}, out <-chan struct{}) {
	if atomic.CompareAndSwapInt32(&qs.notify, 0, 1) {
		go func() {
			ch <- struct{}{}
			<-out
			qs.notify = 0
		}()
	}
}

func (qs *QueueStack[T]) Len() int {
	return int(atomic.LoadInt32(&qs.count))
}

func (qs *QueueStack[T]) loadHead() *Node[T] {
	return (*Node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&qs.head))))
}

func (qs *QueueStack[T]) loadTail() *Node[T] {
	return (*Node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&qs.tail))))
}

func (qs *QueueStack[T]) loadNext(node *Node[T]) *Node[T] {
	return (*Node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&node.next))))
}
