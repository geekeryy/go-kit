package xemail_test

import (
	"context"
	"os"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xemail"
)

func TestSendMail(t *testing.T) {
	c := xconfig.New(context.TODO(), apollo.NewSource("http://apollo.dev.jiangyang.me", "go-kit", "default", os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT")), nil)
	confStr := c.GetString("email")
	err := xemail.New(confStr).SendMail([]string{"1126254578@qq.com"}, "subject", "你好")
	if err != nil {
		t.Error(err)
	}
}
