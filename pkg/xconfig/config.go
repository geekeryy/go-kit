package xconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/gojsonq/v2"
)

var (
	_ IConfig = (*config)(nil)
)

// IConfig 配置接口
// 功能如下：
// 1.热更新 主动轮询/被动通知
// 2.通过key获取配置
// 3.订阅配置变化
type IConfig interface {
	// ReloadConfig 加载配置 被动通知
	// 从配置源加载到内存
	ReloadConfig() error

	// Scan 将配置扫描到结构体
	Scan(v any)

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

// config 配置类
type config struct {
	ctx           context.Context
	cancel        context.CancelFunc
	value         atomic.Value
	source        Source
	subscriber    sync.Map
	watchInterval time.Duration
	log           *xlog.Logger
}

// New 创建配置类
// storeHandler 描述资源以何种形式存储到Config.value 返回需要存储的值
func New(ctx context.Context, source Source, opts ...Option) IConfig {
	if source == nil {
		panic(errors.New("invalid source"))
	}
	_cfg := &config{
		source:     source,
		log:        xlog.New(),
		subscriber: sync.Map{},
	}

	_cfg.ctx, _cfg.cancel = context.WithCancel(ctx)

	for _, o := range opts {
		o(_cfg)
	}

	if _cfg.watchInterval >= time.Second {
		_cfg.watch()
	}

	if err := _cfg.ReloadConfig(); err != nil {
		panic("config:" + err.Error())
	}
	return _cfg
}

func (c *config) Subscribe(key string, ch chan struct{}) {
	c.subscriber.Store(key, ch)
}

func (c *config) UnSubscribe(key string) {
	c.subscriber.Delete(key)
}

// 发布订阅消息
func (c *config) publish() {
	c.subscriber.Range(func(key, value any) bool {
		go func(key, value any) {
			ch, ok := value.(chan struct{})
			if !ok {
				return
			}
			after := time.After(time.Second * 5)
			select {
			case ch <- struct{}{}:
			case <-after:
				log.Println("超时放弃", key)
			}
		}(key, value)
		return true
	})
}

func (c *config) ReloadConfig() error {
	marshal, err := c.source.GetConfig()
	if err != nil {
		return err
	}
	if marshal == nil && len(marshal) == 0 {
		return nil
	}
	if c.value.CompareAndSwap(nil, marshal) {
		c.log.Info(c.ctx, "配置加载", string(marshal))
		return nil
	}
	if !reflect.DeepEqual(c.value.Load(), marshal) {
		c.log.Info(c.ctx, "配置更新", string(marshal))
		c.value.Store(marshal)
		c.publish()
		return nil
	}

	return nil
}

// watch 创建资源监听 处理资源变化 主动轮训
func (c *config) watch() {
	xsync.NewGroup(xsync.WithUUID("config Watch"), xsync.WithContext(c.ctx)).Go(func(ctx context.Context) error {
		ticker := time.NewTicker(c.watchInterval)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if err := c.ReloadConfig(); err != nil {
					log.Println(err)
				}
			}
		}
	})
}

func (c *config) GetString(key string) string {
	cfg, ok := c.value.Load().([]byte)
	if !ok {
		return ""
	}
	value := gojsonq.New().FromString(string(cfg)).Find(key)
	return fmt.Sprint(value)
}

// Scan 扫描配置到结构体
// 写入value时做了严格限制，只要初始化没有报错，就不会报错
func (c *config) Scan(v any) {
	value, _ := c.value.Load().([]byte)
	_ = json.Unmarshal(value, v)
}

func (c *config) Close() {
	c.cancel()
}
