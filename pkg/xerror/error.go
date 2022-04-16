package xerror

import (
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// New 适用于微服务间的错误传递
// code 4位全局统一错误码（服务间通信）/ 5位用户自定义错误码（具体业务）
// msg 用户错误描述，为空时使用默认错误描述
// errs 附加错误信息用于日志输出辅助问题查找
// 错误处理的核心在于包裹，一层层的包裹错误既能读取所有错误信息，又能方便找到错误根因；面向业务的错误处理需要全局统一错误码标识，以及错误详情
func New(code code, msg string, errs ...interface{}) error {
	if len(msg) == 0 {
		msg = code.String()
	}
	s := status.New(codes.Code(code.GetCode()), msg)
	if s.Code() == codes.OK {
		return nil
	}
	p := s.Proto()
	if len(errs) > 0 {
		for _, err := range errs {
			anyErr, err := anypb.New(&spb.Status{
				Message: fmt.Sprintf("%+v", err),
			})
			if err != nil {
				continue
			}
			p.Details = append(p.Details, anyErr)
		}
	}
	return status.FromProto(p).Err()
}

// Wrap 向err中包裹错误信息
func Wrap(err error, errs ...interface{}) error {
	s := status.Convert(err)
	if s.Code() == codes.OK {
		return nil
	}

	p := s.Proto()
	if len(errs) > 0 {
		for _, err := range errs {
			anyErr, err := anypb.New(&spb.Status{
				Message: fmt.Sprintf("%+v", err),
			})
			if err != nil {
				continue
			}
			p.Details = append(p.Details, anyErr)
		}
	}
	return status.FromProto(p).Err()
}

// Cause 返回创建err时的根因Code
func Cause(err error) Code {
	return Code(status.Convert(err).Code())
}
