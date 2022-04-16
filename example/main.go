// Package client @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/4/4 4:22 PM
package main

import (
	"context"
	"log"
	"time"
)

func main() {
	ch := make(chan string)
	time.AfterFunc(time.Second, func() {
		close(ch)
	})
	log.Println(<-ch)

}

func context2() {
	var ch chan struct{} = make(chan struct{})
	close(ch)
	select {
	case <-ch:
		log.Println("ch")
	default:
		log.Println(ch)
	}
}

func context1() {
	ctx, _ := context.WithCancel(context.Background())
	ctx1, _ := context.WithCancel(ctx)
	ctx2 := &MyCtx{ctx1}
	_, _ = context.WithCancel(ctx2)
	<-make(chan struct{})
}

type MyCtx struct{ context.Context }

var done = make(<-chan struct{})

func (*MyCtx) Done() <-chan struct{} {
	return done
}
