package xenv

import (
	"os"
	"strings"
)

var _envMap = make(map[string]string)

const (
	AppEnv       = "APP_ENV"
	AppName      = "APP_NAME"
	AppVersion   = "APP_VERSION"
	TraceName    = "TRACE_NAME"
	ApolloUrl    = "APOLLO_URL"
	ApolloAppID  = "APOLLO_APPID"
	ApolloSecret = "APOLLO_ACCESS_KEY_SECRET"
)

func Init(userConf map[string]string) {
	if appEnv := os.Getenv(AppEnv); len(appEnv) > 0 {
		_envMap[AppEnv] = appEnv
	} else {
		panic("invalid APP_ENV")
	}
	appName := userConf[AppName]
	if len(appName) == 0 {
		if appName = os.Getenv(AppName); len(appName) > 0 {
			_envMap[AppName] = appName
		} else {
			panic("invalid APP_NAME")
		}
	}
	for k, v := range userConf {
		if value := os.Getenv(k); len(value) > 0 {
			_envMap[k] = value
			continue
		}
		if value := os.Getenv(k + "_" + strings.ToUpper(strings.ReplaceAll(appName, "-", "_"))); len(value) > 0 {
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

func GetApolloCluster(cluster string) string {
	switch _envMap[AppEnv] {
	case _local:
		return "local"
	case _test:
		return "test"
	case _dev:
		return "default"
	case _prod:
		return "prod"
	}
	return cluster
}

func GetApolloNamespace(namespace string) string {
	switch _envMap[AppEnv] {
	case _local:
		return "dev." + namespace + ".local"
	case _test:
		return "dev." + namespace + ".test"
	case _dev:
		return "dev." + namespace
	case _prod:
		return "dev." + namespace + ".prod"
	}
	return namespace
}

func GetApolloSecret() string {
	secret := os.Getenv(ApolloSecret)
	switch _envMap[AppEnv] {
	case _local:
		if value := os.Getenv(ApolloSecret + "_" + strings.ToUpper(strings.ReplaceAll(_envMap[AppName], "-", "_"))); len(value) > 0 {
			return value
		}
	case _test:
	case _dev:
	case _prod:
	}
	return secret
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
