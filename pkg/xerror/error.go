package xerror

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type code interface {
	String() string
	GetCode() uint32
}

// New
// code 4位全局统一错误码（服务间通信）/ 5位用户自定义错误码（具体业务）
// msg 用户错误信息，为空时使用默认错误信息
// errs 附加错误信息用于日志输出辅助问题查找
func New(code code, msg string, errs ...interface{}) error {
	if len(msg) == 0 {
		msg = code.String()
	}
	s := status.New(codes.Code(code.GetCode()), msg)
	if s.Code() == codes.OK {
		return nil
	}
	if len(errs) > 0 {
		details := make([]proto.Message, 0)
		for _, err := range errs {
			details = append(details, &spb.Status{
				Message: fmt.Sprintf("%+v", err),
			})
		}
		s, _ = s.WithDetails(details...)
	}
	return s.Err()
}
