package xhttp_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/comeonjy/go-kit/pkg/xhttp"
	"github.com/comeonjy/go-kit/pkg/xsync"
)

func TestHttpClient_Get(t *testing.T) {
	c := xhttp.NewHttp()
	g := xsync.NewGroup()
	log.Println("start")
	for i := 0; i < 100; i++ {
		g.Go(func(ctx context.Context) error {
			for i := 0; i < 1000; i++ {
				statusCode, _, err := c.Get("http://localhost:8080/v1/ping")
				if err != nil {
					t.Error(err)
					//return errors.WithMessage(err, time.Now().Format("2006-01-02 15:04:05"))
				}
				if statusCode != http.StatusOK {
					t.Error(statusCode)
					//return errors.WithMessage(errors.New(fmt.Sprintf("%d", statusCode)), time.Now().Format("2006-01-02 15:04:05"))
				}
			}
			return nil
		})
	}
	g.Wait()

}
func Get(urlStr string) (int, []byte, error) {
	body := bytes.NewBuffer([]byte{})
	req, err := http.NewRequest(http.MethodGet, urlStr, body)
	if err != nil {
		return 0, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
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
