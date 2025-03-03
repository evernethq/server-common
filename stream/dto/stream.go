package dto

import (
	"time"

	"github.com/evernethq/server-common/util/pubsub"
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

// DaemonStreamData 定义我们自己需要的数据结构，只包含我们需要的字段
type DaemonStreamData struct {
	RequestId string
	Method    string
	Payload   string
}

type DaemonStreamDataInterface interface {
	GetRequestId() string
	GetMethod() string
	GetPayload() string
}

type ReplyLog struct {
	Arg       string            // 接收到的参数
	StartTime time.Time         // 接收数据时间
	SendData  *DaemonStreamData // 发送的数据
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
