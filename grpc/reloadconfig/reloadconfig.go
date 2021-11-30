package reloadconfig

import (
	"context"

	"github.com/comeonjy/go-kit/pkg/xconfig"
)

type Server struct {
	conf xconfig.ReloadConfigInterface
}

func NewServer(conf xconfig.ReloadConfigInterface) *Server {
	return &Server{
		conf: conf,
	}
}

func (s *Server) ReloadConfig(context.Context, *Empty) (*Empty, error) {
	if s.conf != nil {
		return &Empty{}, s.conf.ReloadConfig()
	}
	return &Empty{}, nil
}
