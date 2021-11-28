package xemail_test

import (
	"os"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xemail"
)


func TestSendMail(t *testing.T) {
	c := xconfig.New(
		xconfig.WithSource(apollo.NewSource("http://apollo.dev.jiangyang.me", "go-kit", "default", "application", os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT"))),
	)
	confStr := c.Get("email")
	err := xemail.New(confStr).SendMail([]string{"1126254578@qq.com"}, "subject", "你好")
	if err != nil {
		t.Error(err)
	}
}
