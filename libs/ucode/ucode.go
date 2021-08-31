package ucode

import (
	"math/rand"
	"time"
)

var baseStrLists = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
var base = []byte(baseStrLists)

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

var letterRunes = []rune("abcdef#ghijk_lmnop&qrs@tuvwxyzABCD%EFGHIJKLMNOPQRSTUVWXYZ*")

func RandStringRunesNoNum(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
