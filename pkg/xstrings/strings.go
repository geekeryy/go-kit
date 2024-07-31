// Package xstrings @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/7/31 下午2:41
package xstrings

import "math/rand"

func GenerateRandomID() string {
	letters := "ABCDEFGHIJKLMNPQRSTUVWXYZ"
	numbers := "123456789"

	var result string
	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			result += string(letters[rand.Intn(len(letters))])
		} else {
			result += string(numbers[rand.Intn(len(numbers))])
		}
	}

	return result
}
