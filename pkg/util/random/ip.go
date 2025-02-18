package random

import (
	"math/rand"
	"net"
	"time"
)

// GenerateOperatorNATIP 随机生成运营商 NAT ip 地址
func GenerateOperatorNATIP() net.IP {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 固定前10位为100.64.0.0/10
	// 100.64.0.0 的二进制表示为 01100100.01000000.00000000.00000000
	// /10表示前10位是固定的
	// 所以我们需要在余下的 22 位中随机生成
	ip := make(net.IP, 4)
	ip[0] = 100
	ip[1] = 64 + byte(r.Intn(64)) // 64~127之间
	ip[2] = byte(r.Intn(256))
	ip[3] = byte(r.Intn(256))

	return ip
}
