package xconfig

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/thedevsaddam/gojsonq/v2"
)

// IConfig 配置接口
// 功能如下：
// 1.热更新 主动轮训/被动通知
// 2.通过key获取配置
// 3.订阅配置变化
type IConfig interface {
	// Scan 将资源加载到指定结构体
	Scan(v interface{}) error
	// Watch 创建资源监听 处理资源变化
	Watch(func(c *Config)) error
	// Get 支持json/yaml通过键名获取值 不支持数组
	Get(key string) string
	// Subscribe 用户订阅Watch事件
	Subscribe(key string, ch chan struct{})
	// UnSubscribe 用户取消订阅
	UnSubscribe(key string)
	// Reload 被动触发重载配置
	Reload() error
}

// Source 资源接口
type Source interface {
	// Load 加载配置到资源
	Load() error
	// Value 获取资源中的配置数据
	Value() []byte
	// Reload 重载配置到资源
	Reload() error
	// Watch 主动轮训监听配置
	Watch(interval time.Duration) (chan struct{}, error)
	WithContext(ctx context.Context) Source
}

// Config 配置类
type Config struct {
	ctx        context.Context
	subscriber map[string]chan struct{}
	interval   time.Duration
	source     Source
	watchOnce  sync.Once
}

// New 创建配置类
func New(opts ...Option) IConfig {
	_cfg := &Config{
		ctx:        context.Background(),
		subscriber: make(map[string]chan struct{}),
		interval:   time.Second * 30,
	}
	for _, o := range opts {
		o(_cfg)
	}
	if err := _cfg.source.Load(); err != nil {
		panic("Config:" + err.Error())
	}
	return _cfg
}

func (c *Config) Scan(v interface{}) error {
	return json.Unmarshal(c.source.Value(), v)
}

func (c *Config) Subscribe(key string, ch chan struct{}) {
	c.subscriber[key] = ch
}

func (c *Config) UnSubscribe(key string) {
	delete(c.subscriber, key)
}

// 发布订阅消息
func (c *Config) publish() {
	for k, ch := range c.subscriber {
		after := time.After(time.Second * 1)
		select {
		case ch <- struct{}{}:
		case <-after:
			log.Println("超时放弃", k)
		}
	}
}

func (c *Config) Reload() error {
	if err := c.source.Reload(); err != nil {
		return err
	}
	go c.publish()
	return nil
}

func (c *Config) Watch(handle func(c *Config)) error {
	diff, err := c.source.Watch(c.interval)
	if err != nil {
		return err
	}
	c.watchOnce.Do(func() {
		xsync.NewGroup(xsync.WithUUID("Config Watch"), xsync.WithContext(c.ctx)).Go(func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case _, ok := <-diff:
					if !ok {
						return nil
					}

					if handle != nil {
						handle(c)
					}

					if len(c.subscriber) > 0 {
						c.publish()
					}

				}
			}
		})
	})
	return nil
}

func (c *Config) Get(key string) string {
	cfg := c.source.Value()
	value := gojsonq.New().FromString(string(cfg)).Find(key)
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}
