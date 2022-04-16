package xerror_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xerror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	log.Println(xerror.New(xerror.Code(10001), "自定义错误")) // 自定义错误

	log.Println(xerror.New(xerror.RedisErr, ""))  // 使用默认错误描述
	log.Println(xerror.New(xerror.AuthErr, "失败")) // 覆盖默认错误描述

	err := xerror.New(xerror.SQLErr, "", "告警通知", errors.New("堆栈信息")) // 附加错误信息
	assert(t, xerror.Cause(err), xerror.SQLErr)

	err = xerror.Wrap(err, "add err1")              // 附加错误信息
	err = xerror.Wrap(err, "add err2")              // 附加错误信息
	log.Println(err, status.Convert(err).Details()) // 打印附加错误信息到服务器日志

	assert(t, xerror.Cause(err), xerror.SQLErr) // 断言错误根因
}

func assert(t *testing.T, src, dsc any) {
	if !reflect.DeepEqual(src, dsc) {
		t.Error(t.Name(), "not equal")
		return
	}
	log.Println(src, dsc, "true")
}
