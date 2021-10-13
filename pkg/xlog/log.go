package xlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc/metadata"
)

var _logger *Logger

func init() {
	_logger = New()
}

type Logger struct {
	traceName   string
	feishuUrl   string
	serviceName string
	env         string
	debug       bool
}

func New(opt ...Option) *Logger {
	l := &Logger{}
	for _, o := range opt {
		o(l)
	}
	return l
}

func Info(ctx context.Context, v ...interface{}) {
	_logger.Info(ctx, v...)
}

func Debug(ctx context.Context, v ...interface{}) {
	_logger.Error(ctx, v...)
}

func Error(ctx context.Context, v ...interface{}) {
	_logger.Error(ctx, v...)
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	var traceID string
	if len(l.traceName) > 0 {
		traceID = l.getTraceID(ctx)
	}
	l.Println(traceID, "INFO", v...)
}

func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	if !l.debug {
		return
	}
	var traceID string
	if len(l.traceName) > 0 {
		traceID = l.getTraceID(ctx)
	}
	if len(l.feishuUrl) > 0 {
		if err := l.postFeishu(v); err != nil {
			l.Println(traceID, "SYSTEM", v...)
		}
	}
	l.Println(traceID, "DEBUG", v...)
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	var traceID string
	if len(l.traceName) > 0 {
		traceID = l.getTraceID(ctx)
	}
	if len(l.feishuUrl) > 0 {
		if err := l.postFeishu(v); err != nil {
			l.Println(traceID, "SYSTEM", v...)
		}
	}
	l.Println(traceID, "ERROR", v...)
}

func (l *Logger) Println(traceID string, types string, v ...interface{}) {
	log.Println(traceID, types, v)
}

// Printf 实现 xmysql.Writer
func (l *Logger) Printf(ctx context.Context, fmtStr string, v ...interface{}) {
	var traceID string
	if len(l.traceName) > 0 {
		traceID = l.getTraceID(ctx)
	}
	log.Println(traceID, fmt.Sprintf(fmtStr, v...))
}

func (l *Logger) getTraceID(ctx context.Context) string {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if t := md.Get(l.traceName); len(t) > 0 {
			return t[0]
		}
	}
	return ""
}

type feiShuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (l *Logger) postFeishu(v ...interface{}) error {
	msg := feiShuMsg{MsgType: "text"}
	msg.Content.Text = fmt.Sprintf("%s [ %s ] : %v", l.serviceName, l.env, v)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", l.feishuUrl, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
