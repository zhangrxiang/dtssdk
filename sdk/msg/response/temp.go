package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
)

type TempRequest struct {
	Request *models.ZoneTempNotify
	Device  *models.SetDeviceRequest
	Value   chan *models.ZoneTempNotify
}

func (t *TempRequest) Handle(ctx context.Context, closer tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	notify := content.(*TempRequest)
	select {
	case t.Value <- notify.Request:
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
