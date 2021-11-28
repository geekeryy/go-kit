package xconfig

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/thedevsaddam/gojsonq/v2"
	"gopkg.in/yaml.v3"

	"github.com/comeonjy/go-kit/pkg/xsync"
)

// IConfig 配置类型的抽象
// 功能如下：
// 1.热更新(尚不完善)TODO
// 2.获取配置
type IConfig interface {
	Scan(v interface{}) error
	Watch(func(c *Config)) error
	Get(key string) string
	Subscribe(key string, ch chan struct{})
	UnSubscribe(key string)
}

type Source interface {
	Load() error
	Value() []byte
	Watch(interval time.Duration) (chan struct{}, error)
	WithContext(ctx context.Context) Source
}

type Config struct {
	ctx        context.Context
	subscriber map[string]chan struct{}
	interval   time.Duration
	source     Source
}

type Option func(*Config)

var _cfg *Config

func New(opts ...Option) IConfig {
	var once sync.Once
	once.Do(func() {
		_cfg = &Config{
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
	})
	return _cfg
}

func (c *Config) Scan(v interface{}) error {
	return json.Unmarshal(c.source.Value(), v)
}

// Subscribe 支持多用户订阅
func (c *Config) Subscribe(key string, ch chan struct{}) {
	c.subscriber[key] = ch
}

func (c *Config) UnSubscribe(key string) {
	delete(c.subscriber, key)
}

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

// Watch 多次watch也只会收到一个通知
func (c *Config) Watch(handle func(c *Config)) error {
	diff, err := c.source.Watch(c.interval)
	if err != nil {
		return err
	}
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

func defaultDecoder(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func yamlDecoder(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
