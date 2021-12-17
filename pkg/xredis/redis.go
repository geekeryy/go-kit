package xredis

import (
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type Options struct {
	Addr         string `json:"addr"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	MaxRetries   int    `json:"max_retries"`
	MinIdleConns int    `json:"min_idle_conns"`
}

func New(confStr string) *redis.Client {
	opt := Options{}
	if err := json.Unmarshal([]byte(confStr), &opt); err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:         opt.Addr,
		Username:     opt.Username,
		Password:     opt.Password,
		DB:           opt.DB,
		MaxRetries:   opt.MaxRetries,
		MinIdleConns: opt.MinIdleConns,
	})
}
