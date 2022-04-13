package xconfig_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func storeHandler(data []byte) interface{} {
	config := Conf{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Println(err)
		return nil
	}
	return config
}

func TestDemo(t *testing.T) {
	c := Conf{}
	f1(&c)
	log.Println(c)
}
func f1(c any) {
	err := json.Unmarshal([]byte(`{"mode":"debug"}`), &c)
	if err != nil {
		log.Println(err)
	}
	log.Println(c)
}

func TestConfig_Get(t *testing.T) {
	t.Run("apollo", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		c := xconfig.New(ctx, apollo.NewSource(xenv.GetEnv(xenv.ApolloUrl), _appId, _clusterName, _secret, "dev.grpc"), storeHandler, xconfig.WithWatch(2*time.Second))

		go func() {
			for {
				fmt.Println(c.GetString("mode"))
				fmt.Println(c.GetString("account_grpc"))
				fmt.Println(c.LoadValue().(Conf))
				time.Sleep(time.Second * 3)
			}
		}()
		<-ctx.Done()
		time.Sleep(time.Second)
	})
	t.Run("yaml", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		c := xconfig.New(ctx, file.NewSource("config.yaml"), storeHandler, xconfig.WithWatch(1*time.Second))

		go func() {
			for {
				fmt.Println(c.GetString("mode"))
				fmt.Println(c.GetString("port"))
				time.Sleep(time.Second * 3)
			}
		}()
		<-ctx.Done()
		time.Sleep(time.Second)
	})

}
