package xjwt_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xjwt"
)

func TestCreateToken(t *testing.T) {
	xjwt.Init("1234")
	token_123 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCdXNpbmVzcyI6IntcInV1aWRcIjpcIjFcIixcInJvbGVcIjoyLFwiZXh0ZW5kXCI6XCJcIn0iLCJleHAiOjE2NTAzMDI4MDQuMzUwNDg3LCJpYXQiOjE2NTAyNjY4MDQuMzUwNDg3fQ.gJqsMZmnpF__k1orHsEMQmgHLXz96qphPTfTD6blLKE"
	marshal, _ := json.Marshal(xjwt.Business{
		UUID: "1",
		Role: 2,
	})
	token, err := xjwt.CreateToken(string(marshal), time.Hour*10)
	log.Println(token, err)

	v := xjwt.Business{}
	if err := xjwt.ParseToken(token_123, &v); err != nil {
		t.Errorf("%+v", err)
	}
	log.Println(v)

}
