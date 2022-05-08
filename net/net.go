// Package net @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/5/8 11:30

//go:build !windows
// +build !windows

package net

import (
	"net"
)

// Listen 包装net.Listen
func Listen(network string, address string) (net.Listener, error) {
	return net.Listen(network, address)
}
