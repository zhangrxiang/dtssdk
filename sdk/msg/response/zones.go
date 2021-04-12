package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"sync"
)

type ZonesResp struct {
	ctx     context.Context
	cancel  context.CancelFunc
	Request *models.GetDefenceZoneReply
	Value   chan *models.Zones
	Zones   *models.Zones
	sync.RWMutex
}

func NewZones(ctx context.Context) *ZonesResp {
	ctx, cancel := context.WithCancel(ctx)
	r := &ZonesResp{
		ctx:     ctx,
		cancel:  cancel,
		Request: new(models.GetDefenceZoneReply),
		Value:   make(chan *models.Zones),
	}
	tao.Register(r.MessageNumber(), r.Unmarshaler, r.Handle)
	return r
}

func (r *ZonesResp) Subscribe(call func(zones *models.Zones)) {
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

func (r *ZonesResp) Handle(ctx context.Context, _ tao.WriteCloser) {
	r.Lock()
	defer r.Unlock()
	content := tao.MessageFromContext(ctx)
	src := content.(*ZonesResp).Request
	if !src.Success {
		return
	}
	r.Zones = &models.Zones{}
	r.Zones.Zones = make([]*models.Zone, len(src.Rows))
	for k, v := range src.Rows {
		r.Zones.Zones[k] = &models.Zone{
			Id:        v.ID,
			Name:      v.ZoneName,
			ChannelId: v.ChannelID,
			Start:     v.Start,
			Finish:    v.Finish,
		}
	}
	select {
	case r.Value <- r.Zones:
	default:
	}
}

func (r *ZonesResp) Unmarshaler(data []byte) (tao.Message, error) {
	r.Lock()
	defer r.Unlock()
	err := proto.Unmarshal(data, r.Request)
	return r, err
}

// Serialize serializes Message into bytes.
func (r *ZonesResp) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (r *ZonesResp) MessageNumber() byte {
	return byte(models.MsgID_GetDefenceZoneReplyID)
}
