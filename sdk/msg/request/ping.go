package request

import (
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"time"
)

type PingReq struct{}

func NewPing() *PingReq {
	return &PingReq{}
}

func (r *PingReq) Write(w tao.WriteCloser) error {
	go func() {
		for {
			err := w.Write(r)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Second * 30)
		}
	}()
	return nil
}

// Serialize serializes Message into bytes.
func (r *PingReq) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// MessageNumber returns message type number.
func (r *PingReq) MessageNumber() byte {
	return byte(250)
}
