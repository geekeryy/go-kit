package xlimiter

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type timeSlot struct {
	timestamp time.Time
	count     int64
}

type Limiter struct {
	sync.Mutex
	windows  []*timeSlot
	width    time.Duration
	size     time.Duration
	maxCount int64
}

func NewLimiter(width time.Duration, size time.Duration, maxCount int64) *Limiter {
	return &Limiter{
		Mutex:    sync.Mutex{},
		windows:  make([]*timeSlot, 0),
		width:    width,
		size:     size,
		maxCount: maxCount,
	}
}

func (m *Limiter) SetMaxCount(maxCount int64) {
	atomic.StoreInt64(&m.maxCount, maxCount)
}

func (m *Limiter) Validate() bool {
	m.Lock()
	defer m.Unlock()

	now := time.Now()
	windowsIndex := -1

	// 删除过期时间槽
	for k, v := range m.windows {
		// v.timestamp > now-m.width
		if v.timestamp.Add(m.width).After(now) {
			break
		}
		windowsIndex = k
	}
	if windowsIndex > -1 {
		m.windows = m.windows[windowsIndex+1:]
	}

	// 判断是否超出限制
	var sum int64
	for _, v := range m.windows {
		sum += v.count
	}
	if sum >= m.maxCount {
		return false
	}

	// 写入窗口数组
	// timestamp > now-size
	if len(m.windows) > 0 && m.windows[len(m.windows)-1].timestamp.Add(m.size).After(now) {
		m.windows[len(m.windows)-1].count++
	} else {
		m.windows = append(m.windows, &timeSlot{
			timestamp: now,
			count:     1,
		})
	}

	return true

}

// Validate redis存储
func Validate(ctx context.Context, cli *redis.Client, key string) bool {
	maxCount := 3
	width := 3 * time.Second

	now := time.Now()
	min := strconv.Itoa(int(now.Add(-width).Unix()))
	max := strconv.Itoa(int(now.Unix()))
	member := uuid.New().String()

	// 删除过期记录
	cli.ZRemRangeByScore(ctx, key, "0", min)

	// 判断是否超限
	if count, _ := cli.ZCount(ctx, key, min, max).Uint64(); count >= uint64(maxCount) {
		return false
	}

	// 写入请求记录
	cli.ZAdd(ctx, key, &redis.Z{Score: float64(now.Unix()), Member: member})

	return true
}

// ValidateScript 使用redis脚本
func ValidateScript(ctx context.Context, cli *redis.Client, key string) bool {
	maxCount := "3"
	width := 3 * time.Second

	now := time.Now()
	min := strconv.Itoa(int(now.Add(-width).Unix()))
	max := strconv.Itoa(int(now.Unix()))
	member := uuid.New().String()
	nowStr := strconv.Itoa(int(now.Unix()))

	script := redis.NewScript(`
local key = KEYS[1]
local min = KEYS[2]
local max = KEYS[3]
local member = KEYS[4]
local maxCount = tonumber(KEYS[5])
local now = tonumber(KEYS[6])+0.0
redis.call('zremrangebyscore',key,'0',min)
local count = redis.call('zcount',key,min,max)
if count >= maxCount then
    return false
end
redis.call('zadd',key,now,member)
return true
`)
	b, err := script.Run(ctx, cli, []string{key, min, max, member, maxCount, nowStr}).Bool()
	if err != nil {
		return false
	}
	return b
}
