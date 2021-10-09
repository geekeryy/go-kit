// @Description  TODO
// @Author  	 jiangyang  
// @Created  	 2020/11/11 4:10 下午
package xmongo_test

import (
	"context"
	"github.com/comeonjy/util/config"
	"github.com/comeonjy/util/mongodb"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"testing"
)

func init()  {
	config.LoadConfig()
	mongodb.Init(config.GetConfig().Mongodb)

}
type user struct {
	Id  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
}

func TestGetConn(t *testing.T) {
	u:=user{}
	err:=mongodb.Conn("user").FindOne(context.Background(),bson.M{}).Decode(&u)
	logrus.Info(err,u)
}

func TestGetConn2(t *testing.T) {
	var w sync.WaitGroup
	u:=user{}
	for i:=0;i<500;i++{
		w.Add(1)
		go func() {

			for i := 0; i < 1; i++ {
				if err:=mongodb.Conn("user").FindOne(context.Background(),bson.M{}).Decode(&u);err!=nil{
					t.Fatal(err)
				}
			}
			w.Done()
		}()
	}
	w.Wait()


}

func BenchmarkGetConn(b *testing.B) {
	u:=user{}
	for i := 0; i < b.N*10; i++ {
		if err:=mongodb.Conn("user").FindOne(context.Background(),bson.M{}).Decode(&u);err!=nil{
			b.Fatal(err)
		}
	}
}
