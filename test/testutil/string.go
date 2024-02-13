package testutil

import (
	"math/rand"
	"time"
)

const (
	RandomStringCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStringCharsetLength = len(RandomStringCharset)
)

func RandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = RandomStringCharset[random.Intn(randomStringCharsetLength)]
	}
	return string(result)
}
