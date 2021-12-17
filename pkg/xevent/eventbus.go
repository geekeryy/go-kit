package xevent

import (
	"sync"
)

type EventBus struct {
	data []string
}

var eventBus *EventBus
var eventBusOnce sync.Once

func NewEventBus() *EventBus {
	eventBusOnce.Do(func() {
		eventBus = &EventBus{}
	})
	return eventBus
}

// Register 注册一个事件
func (e *EventBus) Register() {

}

// Sync 同步事件
func (e *EventBus) Sync() {

}

// Async 异步事件
func (e *EventBus) Async() {

}
