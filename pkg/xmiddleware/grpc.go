// Package xmiddleware @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/9/11 4:12 下午
package xmiddleware

import (
	"context"
	"os"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgrpc"
	"go.elastic.co/apm/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/comeonjy/go-kit/pkg/xlog"
)

type validator interface {
	Validate() error
}

func GrpcValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return nil, status.New(codes.InvalidArgument, "参数错误").Err()
		}
	}
	return handler(ctx, req)
}

type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

func GrpcAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
		var err error
		ctx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "未授权")
		}
	}
	return handler(ctx, req)
}

func GrpcRecover(logger *xlog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "method:", info.FullMethod, ", time:", time.Now().Format("2006-01-02 15:04:05:06"), ", err:", time.Now().Format("2006-01-02 15:04:05:06"), ", fatal%s", string(debug.Stack()))
				err = status.Error(codes.Internal, "内部错误")
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GrpcApm(apmURL, appName, appVersion, appEnv string) grpc.UnaryServerInterceptor {
	if len(apmURL) == 0 {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	if err := os.Setenv("ELASTIC_APM_SERVICE_NAME", appName); err != nil {
		panic("ELASTIC_APM_SERVICE_NAME not find")
	}
	if err := os.Setenv("ELASTIC_APM_SERVER_URL", apmURL); err != nil {
		panic("ELASTIC_APM_SERVER_URL not find")
	}
	apm.DefaultTracer.Close()
	if _, err := transport.InitDefault(); err != nil {
		panic(err)
	}
	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:        appName,
		ServiceVersion:     appVersion,
		ServiceEnvironment: appEnv,
	})
	if err != nil {
		panic(err)
	}
	return apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithTracer(tracer), apmgrpc.WithRecovery())
}

func GrpcLogger(traceName string, logger *xlog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var traceID string
		if mdIn, ok := metadata.FromIncomingContext(ctx); ok {
			if t := mdIn.Get(traceName); len(t) > 0 {
				traceID = t[0]
			}
		}

		if len(traceID) == 0 {
			uuids, _ := uuid.NewUUID()
			traceID = uuids.String()
		}

		ctx = metadata.AppendToOutgoingContext(ctx, traceName, traceID)

		logger.Info(ctx, "GRPC", info.FullMethod, req)
		resp, err = handler(ctx, req)
		if err != nil {
			s := status.Convert(err)
			logger.Error(ctx, s.Code(), s.Message(), s.Details())
		}
		return resp, err
	}
}
