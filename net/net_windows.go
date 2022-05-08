// Package net @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/5/8 11:30

//go:build windows
// +build windows

package net

import (
	"context"
	"net"
	"syscall"
)

// Listen 包装net.Listen
func Listen(network string, address string) (net.Listener, error) {
	config := net.ListenConfig{
		Control: ExclusiveAddrUseControl,
	}
	return config.Listen(context.Background(), network, address)
}

// ExclusiveAddrUseControl windows下禁用端口复用
// SO_EXCLUSIVEADDRUSE=(int)(~SO_REUSEADDR)
func ExclusiveAddrUseControl(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		// windows下编译才能通过
		err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, ^syscall.SO_REUSEADDR, 1)
		if err != nil {
			panic(err)
		}
	})
}
