package xetcd_test

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	//ctx:=context.Background()
	//cli:=xetcd.New()
	//get, err := cli.Get(ctx, "a/",clientv3.WithPrefix())
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//log.Println(get.Kvs)
	//
	//lease := clientv3.NewLease(cli)
	//grant, err := lease.Grant(ctx, 10)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//alive, err := lease.KeepAlive(ctx, grant.ID)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//go func() {
	//	for{
	//		log.Println("续租",<-alive)
	//	}
	//}()
	//kv := clientv3.NewKV(cli)
	//put, err := kv.Put(ctx,"a/b/c1", "",clientv3.WithLease(grant.ID) )
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//log.Println("写入成功",put.Header.Revision)
	//for {
	//	get, err := kv.Get(ctx, "a/b/c1",clientv3.WithPrefix())
	//	if err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	log.Println(get.Kvs)
	//	time.Sleep(time.Second)
	//	if get.Count==0{
	//		log.Println("kv 过期")
	//		return
	//	}
	//}

}

func TestDemo(t *testing.T) {
	const zero = 10 / 3
	log.Printf("%T %v", zero, zero)
}
