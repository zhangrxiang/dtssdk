package request

import (
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
)

type DeviceRequest struct {
	Request *models.SetDeviceRequest
}

// Serialize serializes Message into bytes.
func (t *DeviceRequest) Serialize() ([]byte, error) {
	return proto.Marshal(t.Request)
}

// MessageNumber returns message type number.
func (t *DeviceRequest) MessageNumber() byte {
	return byte(models.MsgID_SetDeviceReplyID)
}
