package file

import (
	"context"
	"io/ioutil"
	"log"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/ghodss/yaml"
)

type file struct {
	ctx     context.Context
	once    sync.Once
	name    string
	content atomic.Value
}

const (
	_YAML = "_YAML"
)

var _extMap = map[string]string{
	".yaml": _YAML,
	".yml":  _YAML,
}

func (f *file) Load() error {
	readFile, err := f.load()
	if err != nil {
		return err
	}
	f.content.Store(readFile)
	return nil
}

func (f *file) Value() []byte {
	return f.content.Load().([]byte)
}

func (f *file) WithContext(ctx context.Context) xconfig.Source {
	f.ctx = ctx
	return f
}

func (f *file) load() ([]byte, error) {
	readFile, err := ioutil.ReadFile(f.name)
	if err != nil {
		return nil, err
	}
	ext := path.Ext(f.name)
	switch _extMap[ext] {
	case _YAML:
		readFile, err = yaml.YAMLToJSON(readFile)
		if err != nil {
			return nil, err
		}
	}
	return readFile, nil
}

func (f *file) Watch(interval time.Duration) (chan struct{}, error) {
	var err error
	var diff chan struct{}
	f.once.Do(func() {
		diff = make(chan struct{})
		xsync.NewGroup(xsync.WithUUID("File Watch"), xsync.WithContext(f.ctx)).Go(func(ctx context.Context) error {
			defer close(diff)
			ticker := time.NewTicker(interval)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					readFile, err := f.load()
					if err != nil {
						log.Println("Config", err)
						continue
					}
					if string(readFile) != string(f.Value()) {
						f.content.Store(readFile)
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
