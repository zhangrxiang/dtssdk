package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"time"
)

type TempRequest struct {
	ctx       context.Context
	cancel    context.CancelFunc
	Request   *models.ZoneTempNotify
	Value     chan *models.ZonesTemp
	ZonesTemp *models.ZonesTemp
}

func NewTempRequest(ctx context.Context) *TempRequest {
	ctx, cancel := context.WithCancel(ctx)
	t := &TempRequest{
		ctx:     ctx,
		cancel:  cancel,
		Request: new(models.ZoneTempNotify),
		Value:   make(chan *models.ZonesTemp, 100),
	}
	tao.Register(t.MessageNumber(), t.Unmarshaler, t.Handle)
	return t
}

func (t *TempRequest) Subscribe(call func(*models.ZonesTemp)) {
	go func() {
		ticket := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticket.C:
				log.Println("get temp timeout 1 minute,so break")
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

func (t *TempRequest) Handle(ctx context.Context, _ tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	r := content.(*TempRequest).Request
	if t.ZonesTemp == nil || len(t.ZonesTemp.Zones) != len(r.Zones) {
		t.ZonesTemp = &models.ZonesTemp{
			DeviceId:  r.DeviceID,
			CreatedAt: time.Unix(r.Timestamp/1000, 0).Format("2006-01-02 15:04:05"),
			Zones:     make([]models.ZoneTemp, len(r.Zones)),
		}
	}
	for k, z := range r.Zones {
		t.ZonesTemp.Zones[k] = models.ZoneTemp{
			Zone: models.Zone{
				Id:   z.ID,
				Name: z.ZoneName,
			},
			Temperature: models.Temperature{
				Max: z.MaxTemperature,
				Avg: z.AverageTemperature,
				Min: z.MinTemperature,
			},
		}
	}
	select {
	case t.Value <- t.ZonesTemp:
	default:
	}
}

func (t *TempRequest) Unmarshaler(data []byte) (tao.Message, error) {
	err := proto.Unmarshal(data, t.Request)
	return t, err
}

// Serialize serializes Message into bytes.
func (t *TempRequest) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (t *TempRequest) MessageNumber() byte {
	return byte(models.MsgID_ZoneTempNotifyID)
}
