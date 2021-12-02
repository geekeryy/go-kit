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
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
)

type apollo struct {
	url          string
	appId        string
	clusterName  string
	nameSpaceMap map[string]string
	secret       string
	configMap    map[string]interface{}
	mutex        sync.Mutex
}

func NewSource(url string, appId string, clusterName string, secret string, nameSpace ...string) xconfig.Source {
	nameSpaceMap := map[string]string{
		"application": "",
	}
	for _, v := range nameSpace {
		nameSpaceMap[v] = ""
	}
	return &apollo{
		url:          url,
		appId:        appId,
		clusterName:  clusterName,
		nameSpaceMap: nameSpaceMap,
		secret:       secret,
	}
}

func (a *apollo) GetConfig() ([]byte, error) {
	if err := a.getConfig(); err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(a.configMap)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (a *apollo) getConfig() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for k, v := range a.nameSpaceMap {
		aConfigs, err := a.load(k, v)
		log.Println(k, v, aConfigs)
		if err != nil {
			return err
		}
		if aConfigs == nil {
			continue
		}
		a.nameSpaceMap[k] = aConfigs.ReleaseKey
		marshal, err := json.Marshal(aConfigs.Configurations)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(marshal, &a.configMap); err != nil {
			return err
		}
	}
	return nil
}

func (a *apollo) load(nameSpace, releaseKey string) (*apolloConfigs, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	urlStr := fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s", a.url, a.appId, a.clusterName, nameSpace, releaseKey)
	log.Println("POST", urlStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
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
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
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
