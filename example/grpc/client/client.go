// Package client @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/4/4 4:22 PM
package main

import (
	"log"
	"net/http"
)

func main() {
	get, err := http.Get("http://localhost:8080/demo")
	if err != nil {
		log.Println(err)
		return
	}
	cookies := get.Cookies()
	log.Println(cookies)
}
