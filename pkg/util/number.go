package util

import (
	"math/rand"
)

// StrategicGrowthNumber 获取策略增长数
// 计算指定周期内随机增长后的值
// 周期数num 随机区间start end
// 取随机区间中间数为基数，取区间中间一半为实际区间产生随机数 [ [x|x] ]
// 返回 基数*周期数+随机数
// 原理 区间随机数求和总是趋近于区间中间值
func StrategicGrowthNumber(num, start, end int64) int64 {
	var actualStart int64
	if start > end {
		start, end = end, start
	}
	if num < 0 || start < 0 || end <= 0 || start == end {
		return 0
	}

	if end-start < 4 {
		random := rand.Int63n(end-start+1) + start
		return (start+end)*num/2 + random
	}

	base := (start + end) / 2

	if (start+base)%2 == 1 {
		actualStart = (start + base + 1) / 2
	} else {
		actualStart = (start + base) / 2
	}

	actualEnd := (end + base) / 2

	n := actualEnd - actualStart + 1
	random := rand.Int63n(n) + actualStart

	return base*num + random
}
