package xenv_test

import (
	"log"
	"os"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xenv"
)

func TestInit(t *testing.T) {
	var envMap = map[string]string{
		xenv.AppName:     "working",
		xenv.AppVersion:  "v1.0",
		xenv.TraceName:   "trace_id",
		xenv.ApolloAppID: "working",
		xenv.ApolloUrl:   "http://apollo.dev.jiangyang.me",
		"a":              "b",
	}
	os.Setenv(xenv.AppEnv, "local")
	xenv.Init(envMap)
	log.Println(xenv.GetEnv(xenv.AppEnv))
	log.Println(xenv.GetEnv("a"))

}
