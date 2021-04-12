package request

import (
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"sync"
)

type ZonesResp struct {
	channelNum int
	Request    *models.GetDefenceZoneRequest
	sync.RWMutex
}

func NewZones(channelNum int) *ZonesResp {
	r := &ZonesResp{channelNum: channelNum}
	return r
}

func (r *ZonesResp) Write(w tao.WriteCloser) error {
	for i := 1; i <= r.channelNum; i++ {
		r.Request = &models.GetDefenceZoneRequest{
			Channel: int32(i),
			Search:  "",
		}
		err := w.Write(r)
		if err != nil {
			return err
		}
	}
	return nil
}

// Serialize serializes Message into bytes.
func (r *ZonesResp) Serialize() ([]byte, error) {
	return proto.Marshal(r.Request)
}

// MessageNumber returns message type number.
func (r *ZonesResp) MessageNumber() byte {
	return byte(models.MsgID_GetDefenceZoneRequestID)
}
