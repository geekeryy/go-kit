// Package apollo @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/9/5 11:36 上午
package apollo

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xsync"
)

type apollo struct {
	ctx         context.Context
	releaseKey  string
	url         string
	appId       string
	clusterName string
	nameSpace   string
	secret      string
	content     atomic.Value
	once        sync.Once
}

func NewSource(url string, appId string, clusterName string, nameSpace string, secret string) xconfig.Source {
	return &apollo{
		ctx:         context.Background(),
		url:         url,
		appId:       appId,
		clusterName: clusterName,
		nameSpace:   nameSpace,
		secret:      secret,
	}
}

func (a *apollo) WithContext(ctx context.Context) xconfig.Source {
	a.ctx = ctx
	return a
}

func (a *apollo) Load() error {
	aConfigs, err := a.load()
	if err != nil {
		return err
	}
	a.releaseKey = aConfigs.ReleaseKey
	marshal, err := json.Marshal(aConfigs.Configurations)
	if err != nil {
		return err
	}
	a.content.Store(marshal)
	return nil
}

func (a *apollo) Value() []byte {
	return a.content.Load().([]byte)
}

func (a *apollo) Watch() (chan struct{}, error) {
	var diff chan struct{}
	a.once.Do(func() {
		diff = make(chan struct{})
		xsync.NewGroup(xsync.WithContext(a.ctx)).Go(func(ctx context.Context) error {
			defer close(diff)
			// TODO 退避算法
			ticker := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("apollo watcher exit %w", ctx.Err())
				case <-ticker.C:
					get, err := a.load()
					if err != nil {
						log.Println("Config", err)
						continue
					}
					if get.ReleaseKey != a.releaseKey && a.releaseKey != "" {
						a.releaseKey = get.ReleaseKey
						marshal, err := json.Marshal(get.Configurations)
						if err != nil {
							log.Println("Config", err)
							continue
						}
						a.content.Store(marshal)
						diff <- struct{}{}
					}
				}
			}
		})
	})
	return diff, nil
}

func (a *apollo) load() (*apolloConfigs, error) {
	urlStr := fmt.Sprintf("%s/configs/%s/%s/%s", a.url, a.appId, a.clusterName, a.nameSpace)
	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if len(a.secret) > 0 {
		parse, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		timestamp := fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
		sign := signature(timestamp, parse.RequestURI(), a.secret)
		req.Header.Set("Authorization", fmt.Sprintf("Apollo %s:%s", a.appId, sign))
		req.Header.Set("Timestamp", timestamp)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := apolloConfigs{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func signature(timestamp, url, accessKey string) string {
	stringToSign := timestamp + "\n" + url
	key := []byte(accessKey)
	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type apolloConfigs struct {
	AppId          string      `json:"appId"`
	Cluster        string      `json:"cluster"`
	NamespaceName  string      `json:"namespaceName"`
	Configurations interface{} `json:"configurations"`
	ReleaseKey     string      `json:"releaseKey"`
}
