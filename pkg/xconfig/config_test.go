package xconfig_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xconfig/file"
	"github.com/comeonjy/go-kit/pkg/xenv"
)

type Conf struct {
	Mode        string `json:"mode" yaml:"mode"`
	AccountGrpc string `json:"account_grpc"`
	Port        int    `json:"port" yaml:"port"`
}

const (
	_url       = "http://apollo.dev.jiangyang.me"
	_appId     = "go-kit"
	_nameSpace = "application"
)

var (
	_clusterName = "default"
	_secret      = os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT")
)

func TestConfig_Get(t *testing.T) {
	t.Run("apollo", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		c := xconfig.New(ctx, apollo.NewSource(xenv.GetEnv(xenv.ApolloUrl), _appId, _clusterName, _secret, "dev.grpc"), xconfig.WithWatch(2*time.Second))

		go func() {
			for {
				cfg := Conf{}
				c.Scan(&cfg)
				fmt.Println(cfg)
				fmt.Println(c.GetString("mode"))
				fmt.Println(c.GetString("account_grpc"))
				time.Sleep(time.Second * 3)
			}
		}()
		<-ctx.Done()
		time.Sleep(time.Second)
	})
	t.Run("yaml", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		c := xconfig.New(ctx, file.NewSource("config.yaml"), xconfig.WithWatch(1*time.Second))

		go func() {
			for {
				cfg := Conf{}
				c.Scan(&cfg)
				fmt.Println(cfg)
				fmt.Println(c.GetString("mode"))
				fmt.Println(c.GetString("port"))
				time.Sleep(time.Second * 3)
			}
		}()
		<-ctx.Done()
		time.Sleep(time.Second)
	})

}
