package client

type Option struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	ChannelNum int    `json:"channel_num"`
}
