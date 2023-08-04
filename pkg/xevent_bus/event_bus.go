// Package xevent_bus
// 全局事件总线，用于服务间消息通知
package xevent_bus

import (
	"context"
	"github.com/bytedance/godlp/log"
	"github.com/go-redis/redis/v8"
	"go.uber.org/atomic"
	"time"
)

type EventBus struct {
	client *redis.Client
	delay  atomic.Int64 // 能接受的最大延迟时间（单位秒）
}

func (bus *EventBus) WithDelay(delay int64) *EventBus {
	bus.delay.Store(delay)
	return bus
}

func New(redisClient *redis.Client) *EventBus {
	return &EventBus{
		client: redisClient,
	}
}

func (bus *EventBus) Subscribe(ctx context.Context, key string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			pop, err := bus.client.LPop(ctx, key).Result()
			if err != nil {
				if err != redis.Nil {
					log.Errorf("%v", err)
				}
				time.Sleep(bus.delayTime())
				continue
			}
			if len(pop) == 0 {
				log.Errorf("err none pop")
				continue
			}
			ch <- pop
		}
	}()
	return ch
}

func (bus *EventBus) Publish(ctx context.Context, key string, value string) error {
	_, err := bus.client.RPush(ctx, key, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func (bus *EventBus) delayTime() time.Duration {
	delay := bus.delay.Load()
	if delay > 0 {
		return time.Duration(delay) * time.Second
	}
	return 10 * time.Second
}
