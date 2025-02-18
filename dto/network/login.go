package network

import (
	"net"
	"time"

	"github.com/evernethq/server-common/pkg/nebula"
)

// LoginReq 节点登录
type LoginReq struct {
	Token string
	LoginLogSysInfo
}

type LoginLogSysInfo struct {
	Os              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
	Hostname        string
}

type LoginReply struct {
	Name string `json:"name"`
	Type string `json:"type"`
	NebulaConf
}

type NebulaConf struct {
	IP         net.IP       `json:"ip"`
	ExpireTime time.Time    `json:"expire_time"` // 过期时间
	Yaml       *nebula.Conf `json:"yaml"`
}
