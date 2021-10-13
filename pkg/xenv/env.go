package xenv

import (
	"os"
	"strings"
)

var _envMap map[string]string

func Init(appName string, envMap map[string]string) {
	_envMap = make(map[string]string)
	for k, v := range envMap {
		if value := os.Getenv(k); len(value) > 0 {
			_envMap[k] = value
			continue
		}
		if value := os.Getenv(k + "_" + strings.ToUpper(strings.ReplaceAll(appName,"-","_"))); len(value) > 0 {
			_envMap[k] = value
			continue
		}
		_envMap[k] = v
	}
}

func GetEnv(envName string) string {
	if v, ok := _envMap[envName]; ok {
		return v
	}
	return ""
}

// 系统模式
const (
	_modeDebug  = "debug"
	_modeNormal = "normal"
)

// IsDebug 是否为debug模式
func IsDebug(v string) bool {
	return v == _modeDebug
}

// 系统环境枚举
const (
	_prod  = "prod"
	_test  = "test"
	_dev   = "dev"
	_local = "local"
)
