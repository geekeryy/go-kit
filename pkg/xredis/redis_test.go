package xredis_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/go-kit/pkg/xredis"
)

func TestNew(t *testing.T) {
	c := xconfig.New(context.TODO(), apollo.NewSource("http://apollo.dev.jiangyang.me", "go-kit", "default", os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT"), xenv.GetApolloNamespace("grpc"), xenv.GetApolloNamespace("common")), nil)
	cli := xredis.New(c.GetString("redis_conf"))
	arr := make([]string, 0)
	if err := cli.Keys(context.Background(), "*").ScanSlice(&arr); err != nil {
		t.Error(err)
		return
	}
	log.Println(arr)

}
