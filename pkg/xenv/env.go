package xenv

import (
	"os"
	"strings"
)

var _envMap = map[string]string{
	AppEnv:     _local,
	AppName:    "my-app",
	AppVersion: "v1.0.0",
	TraceName:  "trace_id",
	GrpcPort:   "8081",
	HttpPort:   "8080",
	PprofPort:  "6060",
	ApolloUrl:  "http://apollo-service-dev-apollo-configservice.dev:8080",
}

const (
	AppEnv       = "APP_ENV"
	AppName      = "APP_NAME"
	AppVersion   = "APP_VERSION"
	TraceName    = "TRACE_NAME"
	ApolloUrl    = "APOLLO_URL"
	ApolloAppID  = "APOLLO_APPID"
	ApolloSecret = "APOLLO_ACCESS_KEY_SECRET"

	GrpcPort  = "GRPC_PORT"
	HttpPort  = "HTTP_PORT"
	PprofPort = "PPROF_PORT"
)

func init() {
	if appEnv := os.Getenv(AppEnv); len(appEnv) > 0 {
		_envMap[AppEnv] = appEnv
	}
}

// Init 初始化环境配置 ENV > userConf > default
func Init(userConf map[string]string) {
	appName := userConf[AppName]
	if len(appName) > 0 {
		_envMap[AppName] = appName
	} else {
		if appName = os.Getenv(AppName); len(appName) > 0 {
			_envMap[AppName] = appName
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
	if len(secret) == 0 {
		if value := os.Getenv(ApolloSecret + "_" + strings.ToUpper(strings.ReplaceAll(_envMap[AppName], "-", "_"))); len(value) > 0 {
			return value
		}
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

// IsLocal 判断当前环境（环境主要用户切换配置，不可过度依赖）
func IsLocal() bool {
	return _local == _envMap[AppEnv]
}
func IsDev() bool {
	return _dev == _envMap[AppEnv]
}
func IsTest() bool {
	return _test == _envMap[AppEnv]
}
func IsProd() bool {
	return _prod == _envMap[AppEnv]
}

// 系统环境枚举
const (
	_prod  = "prod"
	_test  = "test"
	_dev   = "dev"
	_local = "local"
)
