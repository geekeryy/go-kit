package xlimiter_test

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/go-kit/pkg/xlimiter"
	"github.com/comeonjy/go-kit/pkg/xredis"
	"go.uber.org/atomic"
)

func TestLimiter_Validate(t *testing.T) {
	var w sync.WaitGroup
	var counter atomic.Int64
	limiter := xlimiter.NewLimiter(time.Second*1, time.Microsecond, 10000)
	start := time.Now()
	t.Run("并发", func(t *testing.T) {
		for i := 0; i < 1000000; i++ {
			w.Add(1)
			go func() {
				defer w.Done()
				if limiter.Validate() {
					counter.Add(1)
				}
			}()
		}
		w.Wait()
		log.Println(time.Now().Sub(start), counter)
	})
	t.Run("顺序", func(t *testing.T) {
		for i := 0; i < 100000; i++ {
			w.Add(1)
			if limiter.Validate() {
				counter.Add(1)
			}
			w.Done()
		}
		w.Wait()
		log.Println(time.Now().Sub(start), counter)
	})

}

func BenchmarkLimiter_Validate(b *testing.B) {
	limiter := xlimiter.NewLimiter(time.Minute, time.Second, 10000)
	b.SetParallelism(100000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Validate()
		}
	})
}

func TestValidateScript(t *testing.T) {
	c := xconfig.New(context.TODO(), apollo.NewSource("http://apollo.dev.jiangyang.me", "go-kit", "default", os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT"), xenv.GetApolloNamespace("grpc"), xenv.GetApolloNamespace("common")), nil)
	cli := xredis.New(c.GetString("redis_conf"))
	t.Run("Validate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			log.Println(xlimiter.Validate(context.Background(), cli, "127.0.0.1"))
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 2)
		for i := 0; i < 3; i++ {
			log.Println(xlimiter.Validate(context.Background(), cli, "127.0.0.1"))
			time.Sleep(time.Second)
		}
	})

	t.Run("ValidateScript", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			log.Println(xlimiter.ValidateScript(context.Background(), cli, "127.0.0.2"))
			//time.Sleep(time.Second)
		}
		//time.Sleep(time.Second * 2)
		for i := 0; i < 3; i++ {
			log.Println(xlimiter.ValidateScript(context.Background(), cli, "127.0.0.2"))
			time.Sleep(time.Second)
		}
	})

}
