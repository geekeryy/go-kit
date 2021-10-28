package xconfig

import (
	"context"
	"time"
)

// WithContext 传递上下文
func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		c.ctx = ctx
	}
}

// WithSource 资源配置
func WithSource(s Source) Option {
	return func(c *Config) {
		c.source = s.WithContext(c.ctx)
	}
}

// WithWatchInterval 监控轮训间隙
func WithWatchInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.interval = interval
	}
}
