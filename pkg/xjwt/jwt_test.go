package xjwt_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xjwt"
)

func TestCreateToken(t *testing.T) {
	marshal, _ := json.Marshal(xjwt.Business{
		UUID: "1",
		Role: 2,
	})
	token, err := xjwt.CreateToken(string(marshal), time.Second)
	log.Println(token, err)


	v := xjwt.Business{}
	if err := xjwt.ParseToken(token,&v); err != nil {
		t.Errorf("%+v", err)
	}
	log.Println(v)

}
