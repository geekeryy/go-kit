package xconfig

import (
	"time"
)

type Option func(*Config)

// WithWatch 监控 推荐30秒
func WithWatch(interval time.Duration) Option {
	return func(c *Config) {
		c.watchInterval = interval
	}
}
