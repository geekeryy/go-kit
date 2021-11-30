package xconfig

// ReloadConfigInterface 重载配置Interface
type ReloadConfigInterface interface {
	// ReloadConfig 实现配置重载功能
	ReloadConfig() error
}
