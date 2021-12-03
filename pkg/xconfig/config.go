package xconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/gojsonq/v2"
)

// IConfig 配置接口
// 功能如下：
// 1.热更新 主动轮训/被动通知
// 2.通过key获取配置
// 3.订阅配置变化
type IConfig interface {
	// ReLoad 加载配置 被动通知
	ReLoad() error

	// LoadValue 获取Value
	LoadValue() interface{}

	// GetString 支持json/yaml通过键名获取值 不支持数组
	GetString(key string) string

	// Subscribe 用户订阅Watch事件
	Subscribe(key string, ch chan struct{})
	// UnSubscribe 用户取消订阅
	UnSubscribe(key string)

	// Close 关闭
	Close()
}

// Source 资源接口
type Source interface {
	// GetConfig 获取配置
	GetConfig() ([]byte, error)
}

// Config 配置类
type Config struct {
	ctx           context.Context
	cancel        context.CancelFunc
	value         atomic.Value
	source        Source
	storeHandler  StoreHandlerFun
	subscriber    map[string]chan struct{}
	watchInterval time.Duration
}
type StoreHandlerFun func(data []byte) interface{}

func defaultStoreHandler(data []byte) interface{} {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	return m
}

// New 创建配置类
// storeHandler 描述资源以何种形式存储到Config.value 返回需要存储的值
func New(ctx context.Context, source Source, storeHandler StoreHandlerFun, opts ...Option) IConfig {
	if source == nil {
		panic(errors.New("invalid source"))
	}
	if storeHandler == nil {
		storeHandler = defaultStoreHandler
	}
	_cfg := &Config{
		source:       source,
		storeHandler: storeHandler,
		subscriber:   make(map[string]chan struct{}),
	}

	_cfg.ctx, _cfg.cancel = context.WithCancel(ctx)

	for _, o := range opts {
		o(_cfg)
	}

	if _cfg.watchInterval >= time.Second {
		_cfg.watch()
	}

	if err := _cfg.ReLoad(); err != nil {
		panic("Config:" + err.Error())
	}
	return _cfg
}

func (c *Config) LoadValue() interface{} {
	return c.value.Load()
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

func (c *Config) ReLoad() error {
	marshal, err := c.source.GetConfig()
	if err != nil {
		return err
	}
	if marshal == nil {
		return nil
	}

	value := c.storeHandler(marshal)
	if value == nil {
		return nil
	}
	if c.value.Load() == nil {
		c.value.Store(value)
		return nil
	}
	if !reflect.DeepEqual(c.value.Load(), value) {
		log.Println("配置更新", string(marshal))
		c.value.Store(c.storeHandler(marshal))
		if len(c.subscriber) > 0 {
			go c.publish()
		}
		return nil
	}

	return nil
}

// watch 创建资源监听 处理资源变化 主动轮训
func (c *Config) watch() {
	xsync.NewGroup(xsync.WithUUID("Config Watch"), xsync.WithContext(c.ctx)).Go(func(ctx context.Context) error {
		ticker := time.NewTicker(c.watchInterval)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if err := c.ReLoad(); err != nil {
					log.Println(err)
				}
			}
		}
	})
}

func (c *Config) GetString(key string) string {
	cfg, err := json.Marshal(c.value.Load())
	if err != nil {
		return ""
	}
	value := gojsonq.New().FromString(string(cfg)).Find(key)
	return fmt.Sprint(value)
}

func (c *Config) Close() {
	c.cancel()
}
