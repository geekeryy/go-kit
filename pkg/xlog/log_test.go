package xlog_test

import (
	"context"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"google.golang.org/grpc/metadata"
)

func TestLog(t *testing.T)  {
	logger:=xlog.New(xlog.WithTrace("name"),xlog.WithDebug())
	ctx:=metadata.NewOutgoingContext(context.TODO(),metadata.MD{"name":{"xxxx-xxxx-xxxx-1234"}})
	logger.Info(ctx,"我","是","谁")
	logger.Error(ctx,"我","在","哪")
	logger.Debug(ctx,"D","E","B","U","G")
}
