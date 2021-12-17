// Package xpprof @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/3/23 9:31 PM
package xpprof

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func init() {
	http.DefaultServeMux = http.NewServeMux()

}

func handle(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func Launch(addr string) {
	mux := http.NewServeMux()
	handle(mux)
	go func() {
		log.Fatalln(http.ListenAndServe(addr, mux))
	}()
}
