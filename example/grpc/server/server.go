// Package server @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/4/4 4:22 PM
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/demo", handle)
	log.Println("Running...")
	http.ListenAndServe("localhost:8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Max-Age", "1200")
	log.Println(r.Method, r.URL, r.Cookies())
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: "123",
		//Path:  "/",
		//HttpOnly: true,
	})
	w.Write([]byte(`{"code":0}`))
}
