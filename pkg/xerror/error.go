package xerror

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code interface {
	String() string
	GetCode() uint32
}

// New 输出标准错误信息或者覆盖
func New(code Code, messages ...interface{}) error {
	var msg string
	if len(messages) > 0 {
		msg = fmt.Sprint(messages...)
	} else {
		msg = code.String()
	}
	return status.New(codes.Code(code.GetCode()), msg).Err()
}

// NewError 输出标准错误信息或者覆盖,并包裹附加错误信息用于日志输出辅助问题查找
func NewError(code Code, message string, errs ...interface{}) error {
	msg := code.String()
	if len(message) > 0 {
		msg = message
	}
	s := status.New(codes.Code(code.GetCode()), msg)
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

// DetailsToString 优雅输出detail中的堆栈信息
func DetailsToString(s *status.Status) string {
	var b bytes.Buffer
	for _, v := range s.Details() {
		if errMessage, ok := v.(proto.Message); ok {
			e := spb.Status{}
			errMessageBytes, _ := proto.Marshal(errMessage)
			if err := proto.Unmarshal(errMessageBytes, &e); err == nil {
				str := strings.ReplaceAll(e.GetMessage(), "\t", "")
				arr := strings.Split(str, "\n")
				if len(arr) > 10 {
					arr = arr[:10]
				}
				marshal, _ := json.Marshal(arr)
				b.Write(marshal)
				continue
			}
		}
		b.Write([]byte(fmt.Sprint(v)))
	}

	return b.String()
}

