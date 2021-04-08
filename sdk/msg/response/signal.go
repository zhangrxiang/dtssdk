package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"time"
)

type SignalResp struct {
	ctx           context.Context
	cancel        context.CancelFunc
	Request       *models.TempSignalNotify
	Value         chan *models.ChannelSignal
	ChannelSignal *models.ChannelSignal
}

func NewSignal(ctx context.Context) *SignalResp {
	ctx, cancel := context.WithCancel(ctx)
	s := &SignalResp{
		ctx:     ctx,
		cancel:  cancel,
		Request: new(models.TempSignalNotify),
		Value:   make(chan *models.ChannelSignal, 100),
	}
	tao.Register(s.MessageNumber(), s.Unmarshaler, s.Handle)
	return s
}

func (t *SignalResp) Subscribe(call func(*models.ChannelSignal)) {
	go func() {
		ticket := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticket.C:
				log.Println("get signal timeout 1 minute,so break")
				return
			case t := <-t.Value:
				call(t)
				ticket.Reset(time.Minute)
			case <-t.ctx.Done():
				log.Println("cancel...", t.ctx.Err())
				return
			}
		}
	}()
}

func (t *SignalResp) Handle(ctx context.Context, _ tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	r := content.(*SignalResp).Request
	t.ChannelSignal = &models.ChannelSignal{
		DeviceId:   r.DeviceID,
		ChannelId:  r.ChannelID,
		RealLength: r.RealLength,
		Host:       "",
		Signal:     r.Signal,
		CreatedAt:  time.Unix(r.Timestamp/1000, 0).Format(LocalTimeFormat),
	}
	select {
	case t.Value <- t.ChannelSignal:
	default:
	}
}

func (t *SignalResp) Unmarshaler(data []byte) (tao.Message, error) {
	err := proto.Unmarshal(data, t.Request)
	return t, err
}

// Serialize serializes Message into bytes.
func (t *SignalResp) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (t *SignalResp) MessageNumber() byte {
	return byte(models.MsgID_TempSignalNotifyID)
}
