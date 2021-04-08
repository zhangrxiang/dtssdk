package request

type PingReq struct{}

func NewPing() *PingReq {
	return &PingReq{}
}

// Serialize serializes Message into bytes.
func (t *PingReq) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (t *PingReq) MessageNumber() byte {
	return byte(250)
}
