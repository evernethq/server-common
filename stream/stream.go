package stream

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	pkgErrReason "github.com/evernethq/server-common/errors"
	"github.com/evernethq/server-common/stream/dto"
	"github.com/evernethq/server-common/util/pubsub"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	krErr "github.com/go-kratos/kratos/v2/errors"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/pkg/errors"
)

var (
	ErrChannelInterrupt = krErr.New(504, "CHANNEL_INTERRUPT", "channel interrupt")
	ErrNotConnected     = krErr.New(504, "NOT_CONNECTED", "not connected server") // 未连接上服务器
)

type MethodHandlerFunc func(context.Context, []any, func(any) error) (*string, error)

type Base struct {
	SeedCh        chan any
	Request       cmap.ConcurrentMap[string, *dto.ReqContext] // 自己发送的请求
	Log           *log.Helper
	Conn          Conn                                 // 抽象的接口，用于表示连接，实际使用时需要替换为具体的连接实现
	ReqProcessing cmap.ConcurrentMap[string, struct{}] // 正在处理其他端发过来的请求

	ContextStruct []any
	MethodHandler map[string]MethodHandlerFunc
}

type Conn interface {
	Recv() (*dto.DaemonStreamData, error)
	Send(*dto.DaemonStreamData) error
}

func (s *Base) Read(ctx context.Context) error {
	for {
		req, err := s.Conn.Recv()
		if err != nil {
			return errors.WithStack(err)
		}

		s.logRequest(req)

		if reply, ok := s.Request.Get(req.RequestId); ok {
			s.handleReply(req, reply)
		} else {
			s.handleNewRequest(ctx, req)
		}

		select {
		case <-ctx.Done():
			s.Log.Debug("stream read done")
			return nil
		default:
		}
	}
}

// 处理客户端返回的消息
func (s *Base) handleReply(req *dto.DaemonStreamData, reply *dto.ReqContext) {
	if req.Method == dto.ACK {
		reply.ACK = true
		s.Request.Set(req.RequestId, reply)
		s.Log.Debugf("received ACK for request: %s", req.RequestId)
	} else {
		reply.ReplyCh <- req.Payload
		s.Request.Remove(req.RequestId)
		s.logReply(reply)
	}
}

// 接收处理客户端请求的消息
func (s *Base) handleNewRequest(ctx context.Context, req *dto.DaemonStreamData) {
	if req.Method != dto.ACK {
		s.SeedCh <- &dto.ReplyLog{
			Arg:       req.Payload,
			StartTime: time.Now(),
			SendData: &dto.DaemonStreamData{
				RequestId: req.RequestId,
				Method:    dto.ACK,
				Payload:   "",
			},
		}

		if !s.ReqProcessing.Has(req.RequestId) {
			go func() {
				s.ReqProcessing.Set(req.RequestId, struct{}{})
				s.SeedCh <- s.HandleReadData(ctx, req)
				s.ReqProcessing.Remove(req.RequestId)
			}()
		}
	}
}

func (s *Base) HandleReadData(ctx context.Context, data *dto.DaemonStreamData) *dto.ReplyLog {
	res := &dto.ReplyLog{
		Arg:       data.Payload,
		StartTime: time.Now(),
		SendData: &dto.DaemonStreamData{
			RequestId: data.RequestId,
			Method:    data.Method,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 获取处理函数
	handler, ok := s.MethodHandler[data.Method]
	if !ok {
		return nil
	}

	// 调用处理函数
	parsePayloadFunc := ParsePayload(data.Payload)
	reply, err := handler(ctx, s.ContextStruct, parsePayloadFunc)
	if err != nil {
		se := krErr.FromError(err)
		if strings.Contains(se.Message, "context deadline exceeded") {
			se = krErr.FromError(pkgErrReason.TimeOut)
		}
		marshal, err := json.Marshal(se)
		if err != nil {
			return nil
		}
		res.SendData.Payload = string(marshal)
		return res
	}

	if reply != nil {
		res.SendData.Payload = *reply
	}
	return res
}

func (s *Base) Send(ctx context.Context) error {
	for {
		select {
		case data := <-s.SeedCh:
			switch reqData := data.(type) {
			case *dto.ReqData:
				s.sendRequest(reqData)
			case *dto.ReplyLog:
				s.sendReply(reqData)
			}
		case <-ctx.Done():
			s.Log.Debug("stream send done")
			return nil
		}
	}
}

// sendRequest 发送消息给服务端
func (s *Base) sendRequest(reqData *dto.ReqData) {
	requestID := uuid.New().String()
	if reqData != nil {
		if err := s.Conn.Send(&dto.DaemonStreamData{
			RequestId: requestID,
			Method:    reqData.Method,
			Payload:   reqData.Payload,
		}); err != nil {
			s.logError(requestID, reqData, time.Now(), err, "Send")
			return
		}

		s.Request.Set(requestID, &dto.ReqContext{
			ReqData:   *reqData,
			RequestID: requestID,
			ACK:       false,
			StartTime: time.Now(),
		})
	}
}

// 返回信息给客户端
func (s *Base) sendReply(replyLog *dto.ReplyLog) {
	if replyLog != nil {
		if err := s.Conn.Send(replyLog.SendData); err != nil {
			s.logError(replyLog.SendData.RequestId, replyLog.Arg, replyLog.StartTime, err, "Send")
			return
		}

		s.logReply(replyLog)
	}
}

func (s *Base) logError(requestID, arg any, startTime time.Time, err error, op string) {
	s.Log.Errorf("id=%s, args=%v, operation=%s, latency=%v, err=%v", requestID, arg, op, time.Since(startTime), err)
}

// Monitor 监控发送的数据是否接收到返回
func (s *Base) Monitor(ctx context.Context) {
	retryMap := make(map[string]int)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.Log.Debug("stream monitor done")
			return
		case <-ticker.C:
			s.checkRequests(retryMap)
		}
	}
}

func (s *Base) checkRequests(retryMap map[string]int) {
	for k, v := range s.Request.Items() {
		if time.Since(v.StartTime) > 2*time.Second && !v.ACK {
			count := s.updateRetryCount(retryMap, k)
			if count > 3 {
				s.Log.Errorf("retry 3 times, but still failed: %v", v)
				v.ReplyCh <- pkgErrReason.TimeOut
				s.Request.Remove(k)
				delete(retryMap, k)
			} else {
				s.Log.Debugf("restart send data: %v", v)
				s.SeedCh <- v.ReqData
			}
		}
	}
}

// Forwarding 转发其他函数传递过来的消息
func (s *Base) Forwarding(_ context.Context) error {
	return nil
}

func (s *Base) updateRetryCount(retryMap map[string]int, key string) int {
	count := retryMap[key] + 1
	retryMap[key] = count
	return count
}

// ParsePayload 把 Payload 转为对应的结构体
func ParsePayload(payload string) func(any) error {
	return func(v any) error {
		return json.Unmarshal([]byte(payload), v)
	}
}

// GetForwardingResult 获取转发后的结果
func GetForwardingResult(reply, st any) error {
	switch v := reply.(type) {
	case error:
		return v
	case string:
		errData := &dto.ErrData{}
		if err := json.Unmarshal([]byte(v), errData); err != nil {
			return errors.WithStack(err)
		}

		if errData.Code != 0 {
			return krErr.New(errData.Code, errData.Reason, errData.Message)
		}

		if st != nil {
			return json.Unmarshal([]byte(v), st)
		}
	}

	return errors.New("invalid type")
}

func SendData(ctx context.Context, in *dto.SendDataReq) error {
	marshal, err := json.Marshal(in.Data)
	if err != nil {
		return err
	}

	ch := make(chan any, 1)
	defer close(ch)

	if err := in.PS.Publish(in.Topic, &dto.ForwardingReq{
		Method:  in.Method,
		Args:    string(marshal),
		ReplyCh: ch,
	}); err != nil {
		if errors.Is(err, pubsub.ErrTopicNotFound) {
			return ErrNotConnected
		}
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-ch:
		return GetForwardingResult(res, in.ST)
	}
}

func (s *Base) logRequest(req *dto.DaemonStreamData) {
	s.Log.Debugf("read data: id=%s, args=%s, operation=%s", req.RequestId, req.Payload, req.Method)
}

func (s *Base) logReply(reply interface{}) {
	switch v := reply.(type) {
	case *dto.ReqContext:
		s.Log.Infof("Req id=%s, args=%v, operation=%s, latency=%v", v.RequestID, v.ReqData.Payload, v.ReqData.Method, time.Since(v.StartTime))
	case *dto.ReplyLog:
		s.Log.Infof("Reply id=%s, args=%s, return=%v, operation=%s, latency=%v",
			v.SendData.RequestId, v.Arg, v.SendData.Payload, v.SendData.Method, time.Since(v.StartTime))
	}
}
