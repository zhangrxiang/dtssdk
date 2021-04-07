package response

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
)

type SignResponse struct {
	Request *models.TempSignalNotify
	Device  *models.SetDeviceRequest
	Value   chan *models.TempSignalNotify
}

func (t *SignResponse) Handle(ctx context.Context, closer tao.WriteCloser) {
	content := tao.MessageFromContext(ctx)
	notify := content.(*SignResponse)
	select {
	case t.Value <- notify.Request:
	default:
	}
}

func (t *SignResponse) Unmarshaler(data []byte) (tao.Message, error) {
	err := proto.Unmarshal(data, t.Request)
	return t, err
}

// Serialize serializes Message into bytes.
func (t *SignResponse) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (t *SignResponse) MessageNumber() byte {
	return byte(models.MsgID_TempSignalNotifyID)
}
