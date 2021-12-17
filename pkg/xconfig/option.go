package xconfig

import (
	"time"
)

type Option func(*config)

// WithWatch 监控 推荐30秒
func WithWatch(interval time.Duration) Option {
	return func(c *config) {
		c.watchInterval = interval
	}
}
