package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"time"
)

type EventResp struct {
	ctx          context.Context
	cancel       context.CancelFunc
	Request      *models.DeviceEventNotify
	Value        chan *models.ChannelEvent
	ChannelEvent *models.ChannelEvent
}

func NewEvent(ctx context.Context) *EventResp {
	ctx, cancel := context.WithCancel(ctx)
	r := &EventResp{
		ctx:     ctx,
		cancel:  cancel,
		Request: new(models.DeviceEventNotify),
		Value:   make(chan *models.ChannelEvent),
	}
	tao.Register(r.MessageNumber(), r.Unmarshaler, r.Handle)
	return r
}

func (r *EventResp) Subscribe(call func(*models.ChannelEvent)) {
	go func() {
		for {
			select {
			case v := <-r.Value:
				call(v)
			case <-r.ctx.Done():
				log.Println("cancel...", r.ctx.Err())
				return
			}
		}
	}()
}

func (r *EventResp) Handle(ctx context.Context, _ tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	src := content.(*EventResp).Request
	r.ChannelEvent = &models.ChannelEvent{
		DeviceId:      src.DeviceID,
		ChannelId:     src.ChannelID,
		ChannelLength: src.ChannelLength,
		Host:          "",
		CreatedAt:     time.Unix(src.Timestamp/1000, 0).Format(LocalTimeFormat),
		EventType:     src.EventType,
	}
	select {
	case r.Value <- r.ChannelEvent:
	default:
	}
}

func (r *EventResp) Unmarshaler(data []byte) (tao.Message, error) {
	err := proto.Unmarshal(data, r.Request)
	return r, err
}

// Serialize serializes Message into bytes.
func (r *EventResp) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (r *EventResp) MessageNumber() byte {
	return byte(models.MsgID_DeviceEventNotifyID)
}
