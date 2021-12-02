package file

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/ghodss/yaml"
)

type file struct {
	ctx  context.Context
	name string
}

const (
	_YAML = "_YAML"
)

var _extMap = map[string]string{
	".yaml": _YAML,
	".yml":  _YAML,
}

func (f *file) GetConfig() ([]byte, error) {
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

// NewSource 创建文件配置资源 支持json/yaml格式
func NewSource(ctx context.Context, name string) xconfig.Source {
	return &file{
		ctx:  ctx,
		name: name,
	}
}
