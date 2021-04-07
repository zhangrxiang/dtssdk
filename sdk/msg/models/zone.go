package models

type Temperature struct {
	Max float32   `json:"max"`
	Avg float32   `json:"avg"`
	Min float32   `json:"min"`
	At  TimeLocal `json:"at,omitempty"`
}

type AlarmInfo struct {
	Zone
	Temperature
	Location  float32          `json:"location"`
	AlarmAt   TimeLocal        `json:"alarm_at"`
	AlarmType DefenceAreaState `json:"alarm_type"`
}

type ZoneExtend struct {
	Warehouse string `json:"warehouse,omitempty"`
	Group     string `json:"group,omitempty"`
	Row       int    `json:"row,omitempty"`
	Column    int    `json:"column,omitempty"`
	Layer     int    `json:"layer,omitempty"`
}

type Relay struct {
	Tag      string   `json:"tag"`
	Branches []string `json:"branches"`
}

type Zone struct {
	Id        int32   `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	ChannelId int32   `json:"channel_id,omitempty"`
	Host      string  `json:"host,omitempty"`
	Start     float32 `json:"start,omitempty"`
	Finish    float32 `json:"finish,omitempty"`
	Tag       string  `json:"tag,omitempty"`
	Relays    []Relay `json:"relays,omitempty"`
	ZoneExtend
}
