package validate_code

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wq-fotune-backend/pkg/redis"
)

const (
	VCodeKey = "validate:"
)

func CheckCount(key string) (count int, err error) {
	today := time.Now().Format("2006-01-02")
	redisKey := fmt.Sprintf("%s%s:%s", VCodeKey, today, key)
	data, err := redis.CacheGet(redisKey)
	if err == nil {
		data := strings.Split(string(data), ":") //343345:1   验证码:调用次数
		count, _ = strconv.Atoi(data[1])
		if count > 10 {
			return count, errors.New("今日调用次数已达上限！")
		}
	}
	return count, nil
}

func SaveValidateCode(key string, value string, count int, timeout time.Duration) error {
	today := time.Now().Format("2006-01-02")
	redisKey := fmt.Sprintf("%s%s:%s", VCodeKey, today, key)
	value = fmt.Sprintf("%s:%s", value, strconv.Itoa(count))
	return redis.CacheSet(redisKey, value, timeout)
}

func DeleteValidateCode(key string) {
	today := time.Now().Format("2006-01-02")
	redisKey := fmt.Sprintf("%s%s:%s", VCodeKey, today, key)
	redis.CacheDel(redisKey)
}

func GetValidateCode(key string) (string, error) {
	today := time.Now().Format("2006-01-02")
	redisKey := fmt.Sprintf("%s%s:%s", VCodeKey, today, key)
	data, err := redis.CacheGet(redisKey)
	code := strings.Split(string(data), ":")[0] ////343345:1   验证码:调用次数
	return code, err
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	if strings.HasPrefix(sb.String(), "0") {
		return GenValidateCode(width)
	}
	return sb.String()
}
