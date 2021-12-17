// Package app @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/3/23 8:48 PM
package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"google.golang.org/grpc"
)

type App struct {
	ctx  context.Context
	grpc *grpc.Server
	http *http.Server
}

func NewApp(ctx context.Context, http *http.Server, grpc *grpc.Server) *App {
	return &App{
		ctx:  ctx,
		grpc: grpc,
		http: http,
	}
}

func NewAppWithGrpc(ctx context.Context, grpc *grpc.Server) *App {
	return &App{
		ctx:  ctx,
		grpc: grpc,
	}
}

func NewAppWithHttp(ctx context.Context, http *http.Server) *App {
	return &App{
		ctx:  ctx,
		http: http,
	}
}

func (app *App) Run(cancel context.CancelFunc) error {
	g := xsync.NewGroup(xsync.WithContext(app.ctx))
	if app.grpc != nil {
		g.Go(func(ctx context.Context) error {
			return app.runGrpc()
		})
	}
	if app.http != nil {
		g.Go(func(ctx context.Context) error {
			return app.runHttp()
		})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case sig, _ := <-quit:
			log.Println("signal stop ...")
			if app.grpc != nil {
				app.grpc.GracefulStop()
				log.Println("grpc graceful stop")
			}
			if app.http != nil {
				_ = app.http.Shutdown(app.ctx)
				log.Println("http graceful stop")
			}
			cancel()
			return errors.New(fmt.Sprintf("%v", sig))
		}
	}
}

func (app *App) runHttp() error {
	listen, err := net.Listen("tcp", ":"+xenv.GetEnv(xenv.HttpPort))
	if err != nil {
		return err
	}
	log.Printf("http run success in %s \n", listen.Addr().String())
	return app.http.Serve(listen)
}

func (app *App) runGrpc() error {
	listen, err := net.Listen("tcp", ":"+xenv.GetEnv(xenv.GrpcPort))
	if err != nil {
		return err
	}

	log.Printf("grpc run success in %s \n", listen.Addr().String())
	return app.grpc.Serve(listen)
}
