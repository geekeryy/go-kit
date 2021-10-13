package xconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

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
}

type Source interface {
	Load() error
	Value() []byte
	Watch() (chan struct{}, error)
	WithContext(ctx context.Context) Source
}

type Config struct {
	ctx     context.Context
	source  Source
}

type Option func(*Config)

var _cfg *Config

func New(opts ...Option) IConfig {
	var once sync.Once
	once.Do(func() {
		_cfg = &Config{
			ctx:     context.Background(),
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

// Watch 多次watch也只会收到一个通知
func (c *Config) Watch(handle func(c *Config)) error {
	diff, err := c.source.Watch()
	if err != nil {
		return err
	}
	xsync.NewGroup(xsync.WithContext(c.ctx)).Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("config watch exit %w", ctx.Err())
			case _, ok := <-diff:
				if !ok {
					return nil
				}
				handle(c)
			}
		}
	})
	return nil
}

func (c *Config) Get(key string) string {
	cfg := c.source.Value()
	value := gojsonq.New().FromString(string(cfg)).Find(key)
	return value.(string)
}

func defaultDecoder(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func yamlDecoder(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
