package dto

import (
	"time"

	v1 "github.com/evernethq/server-common/api/network/interface/v1"
	"github.com/evernethq/server-common/pkg/util/pubsub"
)

const (
	ACK = "ack" // 确认消息

	Login            = "login"              // 登录
	ModifyNetwork    = "modify_network"     // 修改网络配置
	GetNetworkConfig = "get_network_config" // 获取网络配置
)

type ReqData struct {
	Method  string   // 方法
	Payload string   // 参数
	ReplyCh chan any // 返回的数据
}

type ReqContext struct {
	ReqData
	RequestID string
	ACK       bool
	StartTime time.Time
}

type ReplyLog struct {
	Arg       string               // 接收到的参数
	StartTime time.Time            // 接收数据时间
	SendData  *v1.DaemonStreamData // 发送的数据
}

type ForwardingReq struct {
	Method  string
	Args    string
	ReplyCh chan any
}

type ErrData struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type SendDataReq struct {
	PS     *pubsub.Publisher
	Topic  string
	Method string
	Data   any
	ST     any
}
