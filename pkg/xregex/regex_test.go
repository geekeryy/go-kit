// Package xregex @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/7/21 14:50
package xregex

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegex(t *testing.T) {
	t.Run("email", func(t *testing.T) {
		tests := []struct {
			param  string
			except bool
		}{
			{"111@qq.com", true},
			{"111@qq.com.cn", true},
			{"1111111111111111111111111111@qq.com.cn", true},
			{"111qq.com", false},
			{"111qq.com", false},
		}
		for _, v := range tests {
			require.Equal(t, v.except, rxEmail.MatchString(v.param))
		}
	})
	t.Run("wechatid", func(t *testing.T) {
		tests := []struct {
			param  string
			except bool
		}{
			{"qwqwqwqwqwqwqw", true},
			{"1qwqwqwqwqwqw", false},
		}
		for _, v := range tests {
			require.Equal(t, v.except, rxWeChatID.MatchString(v.param))
		}
	})

}
