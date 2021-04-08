package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"time"
)

type AlarmResp struct {
	ctx        context.Context
	cancel     context.CancelFunc
	Request    *models.ZoneAlarmNotify
	Value      chan *models.ZonesAlarm
	ZonesAlarm *models.ZonesAlarm
}

func NewAlarm(ctx context.Context) *AlarmResp {
	ctx, cancel := context.WithCancel(ctx)
	r := &AlarmResp{
		ctx:     ctx,
		cancel:  cancel,
		Request: new(models.ZoneAlarmNotify),
		Value:   make(chan *models.ZonesAlarm),
	}
	tao.Register(r.MessageNumber(), r.Unmarshaler, r.Handle)
	return r
}

func (r *AlarmResp) Subscribe(call func(*models.ZonesAlarm)) {
	go func() {
		ticket := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticket.C:
				log.Println("get alarm timeout 1 minute,so break")
				return
			case v := <-r.Value:
				call(v)
				ticket.Reset(time.Minute)
			case <-r.ctx.Done():
				log.Println("cancel...", r.ctx.Err())
				return
			}
		}
	}()
}

func (r *AlarmResp) Handle(ctx context.Context, _ tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	src := content.(*AlarmResp).Request
	r.ZonesAlarm = &models.ZonesAlarm{
		DeviceId:  src.DeviceID,
		Host:      "",
		CreatedAt: time.Unix(src.Timestamp/1000, 0).Format(LocalTimeFormat),
	}
	r.ZonesAlarm.Zones = make([]models.ZoneAlarm, len(src.Zones))
	for k, v := range src.Zones {
		r.ZonesAlarm.Zones[k] = models.ZoneAlarm{
			Zone: models.Zone{
				Id:        v.ID,
				Name:      v.ZoneName,
				ChannelId: v.ChannelID,
			},
			Temperature: models.Temperature{
				Max: v.MaxTemperature,
				Avg: v.AverageTemperature,
				Min: v.MinTemperature,
			},
			Location:  v.AlarmLoc,
			AlarmType: v.AlarmType,
		}
	}
	select {
	case r.Value <- r.ZonesAlarm:
	default:
	}
}

func (r *AlarmResp) Unmarshaler(data []byte) (tao.Message, error) {
	err := proto.Unmarshal(data, r.Request)
	return r, err
}

// Serialize serializes Message into bytes.
func (r *AlarmResp) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (r *AlarmResp) MessageNumber() byte {
	return byte(models.MsgID_ZoneAlarmNotifyID)
}
