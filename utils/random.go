package utils

import (
	"errors"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)

}

func RandomStrArr(arr []string, length int) error {
	rand.Seed(time.Now().UnixNano())
	if len(arr) <= 0 {
		return errors.New("the length of the parameter strings should not be less than 0")
	}

	if length <= 0 || len(arr) < length {
		return errors.New("the size of the parameter length illegal")
	}

	for i := len(arr) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		arr[i], arr[num] = arr[num], arr[i]
	}
	return nil
}

func GetRandNumber(num int) int {
	rand.Seed(time.Now().UnixNano())
	if num == 0 {
		return 0
	}
	res := rand.Intn(num)
	if res < 0 {
		return -res
	}
	return res
}

func GetRandNumberInt64(num int64) int64 {
	rand.Seed(time.Now().UnixNano())
	if num == 0 {
		return 0
	}
	res := rand.Int63n(num)
	if res < 0 {
		return -res
	}
	return res
}
