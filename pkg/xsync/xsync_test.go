package xsync_test

import (
	"context"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xsync"
)

func TestGoroutine_Go(t *testing.T) {
	t.Run("g", func(t *testing.T) {
		var g xsync.Group
		t.Run("doPanic", func(t *testing.T) {
			g.Go(doPanic)
			g.Wait()
		})
	})

	t.Run("g with ctx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		g := xsync.NewGroup(xsync.WithContext(ctx))
		t.Run("do", func(t *testing.T) {
			g.Go(do)
			g.Wait()
		})
	})

	t.Run("g with max go num", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		g := xsync.NewGroup(xsync.WithContext(ctx), xsync.WithMaxGoNum(4))
		t.Run("do", func(t *testing.T) {
			for i := 0; i < 4; i++ {
				g.Go(doPanic)
				g.Go(do)
			}
			g.Wait()
		})
	})

}

func BenchmarkGroup_Go(b *testing.B) {
	b.Run("c10k", func(b *testing.B) {
		b.Run("g", func(b *testing.B) {
			g := xsync.NewGroup(xsync.WithMaxGoNum(100))
			for i := 0; i < b.N; i++ {
				g.Go(func(ctx context.Context) error {
					http.Get("http://localhost:8080/v1/user/2")
					return nil
				})
			}
			g.Wait()
		})
		b.Run("go", func(b *testing.B) {
			w := sync.WaitGroup{}
			for i := 0; i < b.N; i++ {
				w.Add(1)
				go func() {
					defer w.Done()
					http.Get("http://localhost:8080/v1/user/1")
				}()
				w.Wait()
			}
		})
		b.Run("serial", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				http.Get("http://localhost:8080/v1/user/1")
			}
		})

	})
}

func doPanic(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Println("wait...")
			time.Sleep(time.Second * 1)
			panic("wait panic")
		}
	}
}

func do(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Println("wait...")
			time.Sleep(time.Second * 1)
		}
	}
}
