package xconfig_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
)

type Conf struct {
	Mode        string `json:"mode" yaml:"mode"`
	AccountGrpc string `json:"account_grpc"`
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

func storeHandler(c *xconfig.Config, data []byte) bool {
	config := Conf{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Println(err)
		return false
	}
	if reflect.DeepEqual(c.LoadValue(), config) {
		return false
	}
	c.StoreValue(config)
	return true
}

func TestConfig_Get(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	c := xconfig.New(ctx, storeHandler, apollo.NewSource( _url, _appId, _clusterName, _secret, "dev.grpc"), xconfig.WithWatch(30*time.Second))

	go func() {
		for  {
			fmt.Println(c.GetValue("mode"))
			fmt.Println(c.GetValue("account_grpc"))
			time.Sleep(time.Second*3)
		}
	}()
	time.Sleep(time.Hour)
}

//func TestConfig(t *testing.T) {
//	t.Run("apollo", func(t *testing.T) {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//		defer cancel()
//		c := xconfig.New(
//			xconfig.WithContext(ctx),
//			xconfig.WithWatchInterval(time.Second*2),
//			xconfig.WithSource(apollo.NewSource(_url, _appId, _clusterName, _nameSpace, _secret)),
//		)
//
//		var vConf atomic.Value
//		var tempConf Conf
//
//		if err := c.Scan(&tempConf); err != nil {
//			t.Error(err)
//			return
//		}
//		vConf.Store(tempConf)
//
//		xsync.NewGroup().Go(func(ctx context.Context) error {
//			for {
//				log.Println("go", vConf.Load().(Conf))
//				time.Sleep(time.Second)
//			}
//		})
//
//		if err := c.Watch(func(config *xconfig.Config) {
//			if err := config.Scan(&tempConf); err != nil {
//				t.Error(err)
//				return
//			}
//			vConf.Store(tempConf)
//		}); err != nil {
//			t.Error(err)
//			return
//		}
//
//		<-ctx.Done()
//	})
//	t.Run("yaml", func(t *testing.T) {
//
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//		defer cancel()
//		c := xconfig.New(
//			xconfig.WithContext(ctx),
//			xconfig.WithWatchInterval(time.Second*2),
//			xconfig.WithSource(file.NewSource("config.yaml")),
//		)
//
//		var vConf atomic.Value
//		var tempConf Conf
//
//		if err := c.Scan(&tempConf); err != nil {
//			t.Error(err)
//			return
//		}
//		vConf.Store(tempConf)
//
//		xsync.NewGroup().Go(func(ctx context.Context) error {
//			for {
//				log.Println("go", vConf.Load().(Conf))
//				time.Sleep(time.Second)
//			}
//		})
//		ch1 := make(chan struct{}, 1)
//		ch2 := make(chan struct{}, 1)
//		c.Subscribe("d1", ch1)
//		c.Subscribe("d2", ch2)
//
//		xsync.NewGroup().Go(func(ctx context.Context) error {
//			for {
//				log.Println("ch1", <-ch1)
//			}
//		})
//		xsync.NewGroup().Go(func(ctx context.Context) error {
//			for {
//				log.Println("ch2", <-ch2)
//			}
//		})
//
//		xsync.NewGroup().Go(func(ctx context.Context) error {
//			for {
//				if err := c.Load(); err != nil {
//					log.Println(err)
//				} else {
//					log.Println("reload")
//				}
//				if err := c.Scan(&tempConf); err != nil {
//					log.Println(err)
//				}
//				vConf.Store(tempConf)
//				time.Sleep(time.Second * 5)
//			}
//		})
//
//		//if err := c.Watch(func(config *xconfig.Config) {
//		//	if err := config.Scan(&tempConf); err != nil {
//		//		t.Error(err)
//		//		return
//		//	}
//		//	vConf.Store(tempConf)
//		//	log.Println("w1")
//		//}); err != nil {
//		//	t.Error(err)
//		//	return
//		//}
//		//
//		//if err := c.Watch(func(config *xconfig.Config) {
//		//	if err := config.Scan(&tempConf); err != nil {
//		//		t.Error(err)
//		//		return
//		//	}
//		//	vConf.Store(tempConf)
//		//	log.Println("w2")
//		//}); err != nil {
//		//	t.Error(err)
//		//	return
//		//}
//
//		<-ctx.Done()
//		time.Sleep(time.Second)
//	})
//}
