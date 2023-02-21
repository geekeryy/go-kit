package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/comeonjy/go-kit/pkg/util"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*100)

	s := &Spider{
		index:    "https://www.aliyun.com/",
		maxLevel: 2,
		savePath: "./html/",
		ctx:      ctx,
	}

	if err := s.Run(); err != nil {
		log.Println(err)
		return
	}
}

// 匹配http链接，取出全路径
var httpReg = regexp.MustCompile("^(http|https)://.*?/(.*)")

// 匹配http链接，取出协议+域名
var domainReg = regexp.MustCompile("^((http|https)://.*?)/")

// 匹配全路径
var relativePathReg = regexp.MustCompile("^/(.*)")

// 匹配style url，取出资源路径
var urlReg = regexp.MustCompile("url\\((.*?)\\)")

// 匹配css url，取出资源路径
var urlCssReg = regexp.MustCompile("url\\(\"(.*?)\"\\)")

type Result struct {
	Sources []Source
}
type Source struct {
	url      string
	filename string
	level    int8
}

type Spider struct {
	index    string
	maxLevel int
	ctx      context.Context
	savePath string
}

func (s *Spider) Run() error {
	if _, err := os.Stat(s.savePath); os.IsNotExist(err) {
		if err := os.Mkdir(s.savePath, os.ModePerm); err != nil {
			return err
		}
	}

	result, err := s.scan(Source{s.index, "index.html", 0})
	if err != nil {
		return err
	}

	for i := 0; i < s.maxLevel; i++ {
		if result == nil || result.Sources == nil {
			return nil
		}
		temp := make([]Source, 0)
		for _, v := range result.Sources {
			if v.level <= int8(s.maxLevel) {
				res, err := s.scan(v)
				if err != nil {
					log.Println(v, err)
					continue
				}
				if res == nil || res.Sources == nil {
					continue
				}
				temp = append(temp, res.Sources...)
			}
		}
		result = &Result{
			Sources: temp,
		}
	}

	return nil

}

// 静态网页
// 1. 下载网页
// 2. 替换所有资源链接
func (s *Spider) scan(source Source) (*Result, error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	default:
	}

	log.Println(source.level, source.url, source.filename)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Path = via[0].URL.Path
			return nil
		},
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(source.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.StatusCode)
	}

	// 处理非html
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// 处理css
		if strings.Contains(resp.Header.Get("Content-Type"), "text/css") {
			sources := make([]Source, 0)
			body = urlCssReg.ReplaceAllFunc(body, func(s []byte) []byte {
				var newStr string
				if backgroundImageSubmatch := urlCssReg.FindSubmatch(s); len(backgroundImageSubmatch) == 2 {
					newStr = string(backgroundImageSubmatch[1])
				} else {
					return s
				}
				suffix := filepath.Ext(strings.SplitN(newStr, "?", 2)[0])
				if len(suffix) == 0 {
					suffix = ".png"
				}
				htmlSubmatch := httpReg.FindStringSubmatch(newStr)
				relativeSubmatch := relativePathReg.FindStringSubmatch(newStr)
				if len(htmlSubmatch) == 3 {
					sources = append(sources, Source{
						url:      newStr,
						filename: util.Md5(newStr) + suffix,
						level:    source.level + 1,
					})
				} else if len(relativeSubmatch) == 2 {
					doaminSubmatch := domainReg.FindStringSubmatch(source.url)
					if len(doaminSubmatch) == 3 {
						sources = append(sources, Source{
							url:      doaminSubmatch[1] + newStr,
							filename: util.Md5(newStr) + suffix,
							level:    source.level + 1,
						})
					} else {
						return s
					}
				} else {
					return s
				}
				return []byte("url(\"/" + util.Md5(newStr) + suffix + "\")")
			})

			body = urlReg.ReplaceAllFunc(body, func(s []byte) []byte {
				var newStr string
				if backgroundImageSubmatch := urlReg.FindSubmatch(s); len(backgroundImageSubmatch) == 2 {
					newStr = string(backgroundImageSubmatch[1])
				} else {
					return s
				}
				suffix := filepath.Ext(strings.SplitN(newStr, "?", 2)[0])
				if len(suffix) == 0 {
					suffix = ".png"
				}
				htmlSubmatch := httpReg.FindStringSubmatch(newStr)
				relativeSubmatch := relativePathReg.FindStringSubmatch(newStr)
				if len(htmlSubmatch) == 3 {
					sources = append(sources, Source{
						url:      newStr,
						filename: util.Md5(newStr) + suffix,
						level:    source.level + 1,
					})
				} else if len(relativeSubmatch) == 2 {
					doaminSubmatch := domainReg.FindStringSubmatch(source.url)
					if len(doaminSubmatch) == 3 {
						sources = append(sources, Source{
							url:      doaminSubmatch[1] + newStr,
							filename: util.Md5(newStr) + suffix,
							level:    source.level + 1,
						})
					} else {
						return s
					}
				} else {
					return s
				}
				return []byte("url(/" + util.Md5(newStr) + suffix + ")")
			})
			// 下载css中的图片字体等资源
			for _, v := range sources {
				if _, err := s.scan(v); err != nil {
					log.Println(v, err)
				}
			}
		}

		return nil, os.WriteFile(s.savePath+source.filename, body, 0666)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	res := Result{}
	urls := make([]Source, 0)

	if source.level < int8(s.maxLevel) {
		doc.Find("a").Each(func(i int, selection *goquery.Selection) {
			urls = append(urls, do(source, "href", ".html", selection)...)
		})
	}

	doc.Find("link").Each(func(i int, selection *goquery.Selection) {
		urls = append(urls, do(source, "href", ".css", selection)...)
	})
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		urls = append(urls, do(source, "src", ".js", selection)...)
	})
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		urls = append(urls, do(source, "src", ".png", selection)...)
	})

	html, err := doc.Html()
	if err != nil {
		return nil, err
	}

	html = urlReg.ReplaceAllStringFunc(html, func(s string) string {
		var newStr string
		if backgroundImageSubmatch := urlReg.FindStringSubmatch(s); len(backgroundImageSubmatch) == 2 {
			newStr = backgroundImageSubmatch[1]
		} else {
			return s
		}
		suffix := filepath.Ext(strings.SplitN(newStr, "?", 2)[0])
		if len(suffix) == 0 {
			suffix = ".png"
		}
		htmlSubmatch := httpReg.FindStringSubmatch(newStr)
		relativeSubmatch := relativePathReg.FindStringSubmatch(newStr)
		if len(htmlSubmatch) == 3 {
			urls = append(urls, Source{
				url:      newStr,
				filename: util.Md5(newStr) + suffix,
				level:    source.level + 1,
			})
		} else if len(relativeSubmatch) == 2 {
			doaminSubmatch := domainReg.FindStringSubmatch(source.url)
			if len(doaminSubmatch) == 3 {
				urls = append(urls, Source{
					url:      doaminSubmatch[1] + newStr,
					filename: util.Md5(newStr) + suffix,
					level:    source.level + 1,
				})
			} else {
				return s
			}
		} else {
			return s
		}
		return "url(/" + util.Md5(newStr) + suffix + ")"
	})

	res.Sources = urls
	return &res, os.WriteFile(s.savePath+source.filename, []byte(html), 0666)

}

func do(source Source, attr, suffix string, selection *goquery.Selection) []Source {
	val, exists := selection.Attr(attr)
	if !exists || len(val) == 0 {
		return nil
	}
	if suffixExt := filepath.Ext(strings.SplitN(val, "?", 2)[0]); len(suffixExt) > 0 {
		suffix = suffixExt
	}
	urls := make([]Source, 0)
	htmlSubmatch := httpReg.FindStringSubmatch(val)
	relativeSubmatch := relativePathReg.FindStringSubmatch(val)
	if len(htmlSubmatch) == 3 {
		selection.SetAttr(attr, "/"+util.Md5(val)+suffix)
		urls = append(urls, Source{
			url:      val,
			filename: util.Md5(val) + suffix,
			level:    source.level + 1,
		})
	}
	if len(relativeSubmatch) == 2 {
		doaminSubmatch := domainReg.FindStringSubmatch(source.url)
		if len(doaminSubmatch) == 3 {
			selection.SetAttr(attr, "/"+util.Md5(val)+suffix)
			urls = append(urls, Source{
				url:      doaminSubmatch[1] + val,
				filename: util.Md5(val) + suffix,
				level:    source.level + 1,
			})
		}
	}
	return urls
}
