package xmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/comeonjy/go-kit/pkg/xerror"
	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func HttpErrorHandler(logger *xlog.Logger) runtime.ErrorHandlerFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		s := status.Convert(err)
		resp := fmt.Sprintf(`{"code":%d,"msg":"%s"}`, int(s.Code()), s.Message())
		logger.Error(ctx, resp, xerror.DetailsToString(s))
		w.Write([]byte(resp))
	}
}

func HttpLogger(traceName string, logger *xlog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get(traceName)
			if len(traceID) == 0 {
				uuids, _ := uuid.NewUUID()
				traceID = uuids.String()
			}
			ctx := metadata.AppendToOutgoingContext(r.Context(), traceName, traceID)
			r = r.WithContext(ctx)
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			logger.Info(ctx, r.Method, r.URL, string(b))
			next.ServeHTTP(w, r)
		})
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func HttpUse(h http.Handler, opt ...MiddlewareFunc) http.Handler {
	for _, v := range opt {
		h = v(h)
	}
	return h
}
