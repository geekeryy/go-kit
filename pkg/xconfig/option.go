package xconfig

import "context"

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