package random

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	Symbols      = "!@#$%^&*()-_+=<>?"
)

const (
	letterIdBits = 6
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

var src = rand.NewSource(time.Now().UnixNano())

func GenerateRandomString(length int, include ...string) string {
	letters := digits
	for _, s := range include {
		letters += s
	}

	b := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
