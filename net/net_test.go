package net_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"syscall"
	"testing"
)

type Result struct {
	Stype   int
	Name    string
	First   string
	Second  string
	Address string
}

var ControlType = "port1"
var Address = ":80"

func TestNet(t *testing.T) {
	tests := []struct {
		Stype   int
		First   string
		Second  string
		Address string
	}{
		{0, "tcp", "tcp", Address},
		{0, "tcp", "tcp4", Address},
		{0, "tcp", "tcp6", Address},
		{0, "tcp4", "tcp", Address},
		{0, "tcp4", "tcp4", Address},
		{0, "tcp4", "tcp6", Address},
		{0, "tcp6", "tcp", Address},
		{0, "tcp6", "tcp4", Address},
		{0, "tcp6", "tcp6", Address},
		{1, "tcp", "tcp", Address},
		{1, "tcp", "tcp4", Address},
		{1, "tcp", "tcp6", Address},
		{1, "tcp4", "tcp", Address},
		{1, "tcp4", "tcp4", Address},
		{1, "tcp4", "tcp6", Address},
		{1, "tcp6", "tcp", Address},
		{1, "tcp6", "tcp4", Address},
		{1, "tcp6", "tcp6", Address},
		{2, "tcp", "tcp", Address},
		{2, "tcp", "tcp4", Address},
		{2, "tcp", "tcp6", Address},
		{2, "tcp4", "tcp", Address},
		{2, "tcp4", "tcp4", Address},
		{2, "tcp4", "tcp6", Address},
		{2, "tcp6", "tcp", Address},
		{2, "tcp6", "tcp4", Address},
		{2, "tcp6", "tcp6", Address},
		{3, "tcp", "tcp", Address},
		{3, "tcp", "tcp4", Address},
		{3, "tcp", "tcp6", Address},
		{3, "tcp4", "tcp", Address},
		{3, "tcp4", "tcp4", Address},
		{3, "tcp4", "tcp6", Address},
		{3, "tcp6", "tcp", Address},
		{3, "tcp6", "tcp4", Address},
		{3, "tcp6", "tcp6", Address},
	}

	result := make([]Result, 0)

	for _, v := range tests {
		name := strconv.Itoa(v.Stype) + "-" + v.First + "-" + v.Second
		t.Run(name, func(t *testing.T) {
			var l1, l2 net.Listener
			switch v.Stype {
			case 0:
				l1, _ = f1(t, v.First, v.Address)
				l2, _ = f1(t, v.Second, v.Address)
			case 1:
				l1, _ = f1(t, v.First, v.Address)
				l2, _ = f2(t, v.Second, v.Address)
			case 2:
				l1, _ = f2(t, v.First, v.Address)
				l2, _ = f1(t, v.Second, v.Address)
			case 3:
				l1, _ = f2(t, v.First, v.Address)
				l2, _ = f2(t, v.Second, v.Address)
			}

			if l1 != nil {
				l1.Close()
			}
			if l2 != nil {
				l2.Close()
			}
			if l1 != nil && l2 != nil {
				result = append(result, Result{
					Stype:   v.Stype,
					Name:    name,
					First:   v.First,
					Second:  v.Second,
					Address: v.Address,
				})
			}
		})
	}

	for _, v := range result {
		var stype string
		switch v.Stype {
		case 0:
			stype = fmt.Sprintf("全net.Listen %s %s", v.First, v.Second)
		case 1:
			stype = fmt.Sprintf("先net.Listen %s 后socketopt:%s %s", v.First, ControlType, v.Second)
		case 2:
			stype = fmt.Sprintf("先socketopt:%s %s 后net.Listen %s", ControlType, v.First, v.Second)
		case 3:
			stype = fmt.Sprintf("全socketopt:%s %s %s", ControlType, v.First, v.Second)
		}
		log.Println(fmt.Sprintf("%s\t%s\t%s", v.Name, stype, "允许复用"))

	}
}

// Port0Control (unix默认)
// 特定地址 不允许任何重用
// 通配符 unix默认情况
func Port0Control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 0)
		if err != nil {
			panic(err)
		}
	})
}

// Port1Control	(unix下使用)
// 特定地址 当双方均开启SO_RESUSEPORT时，允许任意组合（ip需要和协议类型匹配）；否则不允许任何重用
// 通配符 当双方均开启SO_RESUSEPORT时，允许任意组合；否则unix默认情况
func Port1Control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
		if err != nil {
			panic(err)
		}
	})
}

// Addr0Control (windows默认)
// 特定地址 不允许任何重用
// 通配符 默认情况
func Addr0Control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0)
		if err != nil {
			panic(err)
		}
	})
}

// Addr1Control (windows下使用)
// 特定地址 当双方均开启SO_REUSEADDR时，允许任意组合（ip需要和协议类型匹配）；否则不允许任何重用
// 通配符 当双方均开启SO_REUSEADDR时，允许任意组合；否则默认情况
func Addr1Control(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			panic(err)
		}
	})
}

// AddrEControl SO_EXCLUSIVEADDRUSE=^syscall.SO_REUSEADDR (windows下使用)
// 特定地址 不允许任何重用
// 通配符 具体情况请看实验结果
func AddrEControl(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, ^syscall.SO_REUSEADDR, 1)
		if err != nil {
			panic(err)
		}
	})
}

func f1(t *testing.T, network string, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		//t.Error(fmt.Errorf("listen f1 %s %v", network, err))
		return nil, err
	}
	return l, nil
}

func f2(t *testing.T, network string, address string) (net.Listener, error) {
	config := net.ListenConfig{}
	switch ControlType {
	case "port0":
		config.Control = Port0Control
	case "port1":
		config.Control = Port1Control
	case "addr0":
		config.Control = Addr0Control
	case "addr1":
		config.Control = Addr1Control
	case "addre":
		config.Control = AddrEControl
	default:
		// 特定地址 不允许任何重用
		// 通配符 默认情况
	}
	l, err := config.Listen(context.Background(), network, address)
	if err != nil {
		//t.Error(fmt.Errorf("listen f2 %s %v", network, err))
		return nil, err
	}
	return l, nil
}
