package xsync

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Group struct {
	uuid       string
	wait       sync.WaitGroup
	cxt        context.Context
	ch         chan struct{}
	errHandler func(ctx context.Context)
}

type Option func(group *Group)

// WithContext 上下文
func WithContext(ctx context.Context) Option {
	return func(group *Group) {
		group.cxt = ctx
	}
}

// WithErrHandler 错误处理
func WithErrHandler(handler func(context.Context)) Option {
	return func(group *Group) {
		group.errHandler = handler
	}
}

// WithMaxGoNum 限流 至少一个goroutine
// 全局协程管理限流必须带上ctx超时管理，否则会阻塞
func WithMaxGoNum(maxGoNum int64) Option {
	if maxGoNum < 1 {
		maxGoNum = 1
	}
	return func(group *Group) {
		group.ch = make(chan struct{}, maxGoNum)
	}
}

// WithUUID 唯一标识
func WithUUID(uuid string) Option {
	return func(group *Group) {
		group.uuid = uuid
	}
}

// NewGroup 创建Group
func NewGroup(option ...Option) *Group {
	g := Group{
		cxt: context.Background(),
	}
	for _, o := range option {
		o(&g)
	}
	return &g
}

// 添加一个goroutine
func (g *Group) add() {
	g.wait.Add(1)
	if g.ch != nil {
		g.ch <- struct{}{}
	}
}

// 完成一个goroutine
func (g *Group) done() {
	if g.ch != nil {
		<-g.ch
	}
	g.wait.Done()
}

// Go 启动协程并捕获异常 异常退出、panic退出
func (g *Group) Go(f func(context.Context) error) {
	if g.cxt == nil {
		g.cxt = context.Background()
	}
	g.add()
	select {
	case <-g.cxt.Done():
		g.done()
	default:
		go func() {
			defer g.done()
			if err := g.do(f); err != nil {
				if g.errHandler != nil {
					g.errHandler(g.cxt)
				} else {
					log.Printf("Go err: %s %+v", g.uuid, err)
				}
			}
		}()
	}

}

// do 执行函数
func (g *Group) do(f func(context.Context) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("go panic : %s", r)
			return
		}
	}()
	return f(g.cxt)
}

func (g *Group) Wait() {
	g.wait.Wait()
}
