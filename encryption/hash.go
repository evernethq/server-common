package encryption

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// MD5 把字符串转换为 md5  值
func MD5(str string) string {
	has := md5.New()
	has.Write([]byte(str))
	hashBytes := has.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

// SHA256 把字符串转换为 sha256 值
func SHA256(str string) string {
	has := sha256.New()
	has.Write([]byte(str))
	hashBytes := has.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
