package request

import "github.com/zing-dev/dts-sdk/sdk/tao"

type Request interface {
	Write(w tao.WriteCloser) error
	tao.Message
}
