package models

type Temperature struct {
	Max float32    `json:"max"`
	Avg float32    `json:"avg"`
	Min float32    `json:"min"`
	At  *TimeLocal `json:"at,omitempty"`
}

type ZoneAlarm struct {
	Zone
	Temperature
	Location  float32          `json:"location"`
	AlarmAt   TimeLocal        `json:"alarm_at"`
	AlarmType DefenceAreaState `json:"alarm_type"`
}

type ZonesAlarm struct {
	Zones     []ZoneAlarm `json:"zones"`
	DeviceId  string      `json:"device_id"`
	Host      string      `json:"host,omitempty"`
	CreatedAt string      `json:"created_at"`
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

type Zones struct {
	ChannelId int32  `json:"channel_id,omitempty"`
	Host      string `json:"host,omitempty"`
	Zones     []*Zone
}

type ZoneTemp struct {
	Zone
	Temperature
}

type ZonesTemp struct {
	DeviceId  string     `json:"device_id"`
	Host      string     `json:"host,omitempty"`
	CreatedAt string     `json:"created_at"`
	Zones     []ZoneTemp `json:"zones"`
}

type ChannelSignal struct {
	DeviceId   string    `json:"device_id"`
	ChannelId  int32     `json:"channel_id"`
	RealLength float32   `json:"real_length"`
	Host       string    `json:"host,omitempty"`
	Signal     []float32 `json:"signal"`
	CreatedAt  string    `json:"created_at"`
}

type ChannelEvent struct {
	Host          string     `json:"host,omitempty"`
	ChannelId     int32      `json:"channel_id"`
	DeviceId      string     `json:"device_id"`
	EventType     FiberState `json:"event_type"`
	ChannelLength float32    `json:"channel_length"`
	CreatedAt     string     `json:"created_at"`
}
