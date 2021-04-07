package request

import "github.com/zing-dev/dts-sdk/sdk/msg/models"

type ConnectMsg struct {
}

// Serialize serializes Message into bytes.
func (t *ConnectMsg) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (t *ConnectMsg) MessageNumber() byte {
	return byte(models.MsgID_ConnectID)
}
