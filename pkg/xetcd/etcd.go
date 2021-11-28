// Package xetcd @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/11/6 6:19 下午
package xetcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

func New() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	return client
}
