// Package xhttp @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/9/5 11:35 上午
package xhttp

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type httpClient struct {
}

func NewHttp() *httpClient {
	return &httpClient{}
}

func (c *httpClient) Get(urlStr string) (int, []byte, error) {
	body := bytes.NewBuffer([]byte{})
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   100,
		},
	}
	req, err := http.NewRequest(http.MethodGet, urlStr, body)
	if err != nil {
		return 0, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, all, nil
}
