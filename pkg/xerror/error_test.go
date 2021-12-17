package xerror_test

import (
	"log"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xerror"
	"google.golang.org/grpc/status"
)

const (
	myErr xerror.Code = 10001 // 自定义错误
)

func TestNew(t *testing.T) {
	log.Println(xerror.New(myErr, "自定义错误")) // 自定义错误

	log.Println(xerror.New(xerror.RedisErr, ""))   // 使用默认错误信息
	log.Println(xerror.New(xerror.SQLErr, "创建失败")) // 覆盖默认错误信息

	err := xerror.New(xerror.SQLErr, "", "告警通知", "堆栈信息")
	log.Println(err, status.Convert(err).Details()) // 打印附加错误信息到服务器日志
}
