package xslice

import (
	"reflect"
)

// Contains 切片是否包含指定值
func Contains(slice interface{}, value interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}

	sliceLen := s.Len()
	if sliceLen > 0 {
		vKind := reflect.ValueOf(value).Kind()
		if s.Index(0).Kind() == vKind {
			switch vKind {
			case reflect.String:
				return ContainsString(slice.([]string), value.(string))
			case reflect.Int:
				return ContainsInt(slice.([]int), value.(int))
			case reflect.Int64:
				return ContainsInt64(slice.([]int64), value.(int64))
			}
		}

		for i := 0; i < sliceLen; i++ {
			if reflect.DeepEqual(s.Index(i).Interface(), value) {
				return true
			}
		}
	}
	return false
}

func ContainsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsInt64(slice []int64, value int64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
