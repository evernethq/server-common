package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GetSecureRandomElement[T any](slice []T) (T, error) {
	if len(slice) == 0 {
		var zero T
		return zero, fmt.Errorf("empty slice")
	}

	// 生成密码学安全的随机索引
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
	if err != nil {
		var zero T
		return zero, err
	}

	return slice[n.Int64()], nil
}
