package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xsync"
)

type file struct {
	ctx     context.Context
	once    sync.Once
	name    string
	content atomic.Value
}

func (f *file) Value() []byte {
	return f.content.Load().([]byte)
}

func (f *file) Load() error {
	readFile, err := f.load()
	if err != nil {
		return err
	}
	f.content.Store(readFile)
	return nil
}

func (f *file) WithContext(ctx context.Context) xconfig.Source {
	newCtx, _ := context.WithCancel(ctx)
	f.ctx = newCtx
	return f
}

func (f *file) load() ([]byte, error) {
	return ioutil.ReadFile(f.name)
}

func (f *file) Watch() (chan struct{}, error) {
	var err error
	var diff chan struct{}
	f.once.Do(func() {
		diff = make(chan struct{})
		xsync.NewGroup(xsync.WithContext(f.ctx)).Go(func(ctx context.Context) error {
			defer close(diff)
			ticker := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("file watcher exit %w", ctx.Err())
				case <-ticker.C:
					readFile, err := f.load()
					if err != nil {
						log.Println("Config", err)
						continue
					}
					if string(readFile) != string(f.Value()) {
						f.content.Store(readFile)
						log.Println(f.content.Load())
						diff <- struct{}{}
					}
				}
			}
		})
	})
	return diff, err
}

func NewSource(name string) xconfig.Source {
	return &file{name: name}
}
