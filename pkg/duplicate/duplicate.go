// Package duplicate @Description  去重器
// @Author  	 jiangyang
// @Created  	 2023/8/22 11:37
//
// 指定时间内，相同key的函数只执行一次，避免重复提交
package duplicate

import (
	"sync"
	"time"
)

type Deduplicator struct {
	data      map[string]func()
	ttl       time.Duration
	mu        sync.Mutex
	executeCh chan string
}

var globalDedup *Deduplicator

var once sync.Once

func InitDeduplicator(ttl time.Duration) {
	once.Do(func() {
		globalDedup = New(ttl)
	})
}

func Add(key string, f func(), ttl ...time.Duration) {
	globalDedup.Add(key, f, ttl...)
}

func New(ttl time.Duration) *Deduplicator {
	dup := &Deduplicator{
		data:      make(map[string]func()),
		ttl:       ttl,
		executeCh: make(chan string),
	}

	go dup.processor()

	return dup
}

func (d *Deduplicator) Add(key string, f func(), ttl ...time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.data[key]; !exists {
		d.data[key] = f
		ttlTimer := d.ttl
		if len(ttl) > 0 {
			ttlTimer = ttl[0]
		}
		time.AfterFunc(ttlTimer, func() {
			d.executeCh <- key
		})
	}
}

func (d *Deduplicator) processor() {
	for key := range d.executeCh {
		d.mu.Lock()
		if f, exists := d.data[key]; exists {
			f()
			delete(d.data, key)
		}
		d.mu.Unlock()
	}
}
