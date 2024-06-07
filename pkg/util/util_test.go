package util_test

import (
	"github.com/comeonjy/go-kit/pkg/util"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestMax(t *testing.T) {
	type Param struct {
		name  string
		param any
		want  any
	}

	tests := []Param{
		{"整数", []int{1, 2, 3, 4}, 4},
		{"字符串整数", []string{"1", "2", "3", "4", "10"}, "4"},
		{"等长字符串", []string{"zz", "aa", "bb", "cc"}, "zz"},
		{"不等长字符串", []string{"abc", "abb", "a", "ab"}, "abc"},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			switch tt.param.(type) {
			case []int:
				assert.Equal(t, util.Max(tt.param.([]int)), tt.want)
			case []string:
				assert.Equal(t, util.Max(tt.param.([]string)), tt.want)
			}

		})
	}
}
