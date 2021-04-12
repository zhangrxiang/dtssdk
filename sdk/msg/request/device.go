package request

import (
	"github.com/golang/protobuf/proto"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/tao"
)

type DeviceReq struct {
	Request *models.SetDeviceRequest
}

func NewDevice() *DeviceReq {
	return &DeviceReq{
		Request: &models.SetDeviceRequest{
			ZoneTempNotifyEnable:    true,
			ZoneAlarmNotifyEnable:   true,
			FiberStatusNotifyEnable: true,
			TempSignalNotifyEnable:  true,
		}}
}

func (r *DeviceReq) Write(w tao.WriteCloser) error {
	return w.Write(r)
}

// Serialize serializes Message into bytes.
func (r *DeviceReq) Serialize() ([]byte, error) {
	return proto.Marshal(r.Request)
}

// MessageNumber returns message type number.
func (r *DeviceReq) MessageNumber() byte {
	return byte(models.MsgID_SetDeviceReplyID)
}
