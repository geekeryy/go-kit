package xredis_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xredis"
)

func TestNew(t *testing.T) {
	opt := xredis.Options{
		Addr:     "redis.middleware.svc.cluster.local:6379",
		Username: "",
		Password: "",
		DB:       0,
	}
	marshal, err := json.Marshal(opt)
	if err != nil {
		return
	}
	cli := xredis.New(string(marshal))
	arr := make([]string, 0)
	if err := cli.Keys(context.Background(), "*").ScanSlice(&arr); err != nil {
		t.Error(err)
		return
	}
	log.Println(arr)

}
