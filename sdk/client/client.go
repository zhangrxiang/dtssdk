package client

import (
	"context"
	"fmt"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/msg/request"
	"github.com/zing-dev/dts-sdk/sdk/msg/response"
	"github.com/zing-dev/dts-sdk/sdk/tao"
	"log"
	"net"
	"sync"
)

type TopicType byte

const (
	TopicTemp TopicType = iota
	TopicSignal
	TopicAlarm
	TopicEvent
	TopicZones
	TopicDeviceId
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	Option
	conn   *tao.ClientConn
	group  *sync.WaitGroup
	zones  []*models.Zone
	locker sync.RWMutex

	requests map[byte]request.Request
}

func New(o Option) *App {
	ctx, cancel := context.WithCancel(context.Background())
	a := &App{
		ctx:      ctx,
		cancel:   cancel,
		Option:   o,
		group:    new(sync.WaitGroup),
		requests: map[byte]request.Request{},
	}
	a.Publish(request.NewPing())
	a.Publish(request.NewDevice())
	return a
}

func (a *App) valid() {
	if a.Option.Ip == "" {
		panic("ip")
	}

	if a.Option.Port < 1024 {
		panic("port")
	}
}

func (a *App) GetZones() []*models.Zone {
	a.locker.RLock()
	defer a.locker.RUnlock()
	return a.zones
}

func (a *App) Run() {
	a.valid()
	fmt.Println("start", fmt.Sprintf("%s:%d", a.Ip, a.Port))
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.Ip, a.Port))
	if err != nil {
		log.Fatalln(err)
	}
	onConnect := tao.OnConnectOption(func(conn tao.WriteCloser) bool {
		for _, r := range a.requests {
			_ = r.Write(conn)
		}
		return true
	})

	onError := tao.OnErrorOption(func(conn tao.WriteCloser) {
		log.Println("on error")
	})

	onClose := tao.OnCloseOption(func(conn tao.WriteCloser) {
		log.Println("on close")
	})

	onMessage := tao.OnMessageOption(func(msg tao.Message, conn tao.WriteCloser) {
		log.Println(msg.MessageNumber())
	})
	a.conn = tao.NewClientConn(0, c, onConnect, onError, onClose, onMessage, tao.ReconnectOption())
	a.conn.Start()
	a.group.Wait()
}

func (a *App) Publish(requests ...request.Request) {
	a.locker.Lock()
	defer a.locker.Unlock()
	for _, r := range requests {
		a.requests[r.MessageNumber()] = r
		tao.Register(r.MessageNumber(), nil, nil)
	}
}

func (a *App) Subscribe(topic TopicType, call func(result interface{})) {
	switch topic {
	case TopicZones:
		log.Println("zones start")
		response.NewZones(a.ctx).Subscribe(func(zones *models.Zones) {
			a.zones = append(a.zones, zones.Zones...)
			call(zones)
		})
	case TopicTemp:
		response.NewTemp(a.ctx).Subscribe(func(temp *models.ZonesTemp) {
			call(temp)
		})
	case TopicSignal:
		response.NewSignal(a.ctx).Subscribe(func(signal *models.ChannelSignal) {
			call(signal)
		})
	case TopicAlarm:
		response.NewAlarm(a.ctx).Subscribe(func(alarms *models.ZonesAlarm) {
			call(alarms)
		})
	case TopicEvent:
		response.NewEvent(a.ctx).Subscribe(func(e *models.ChannelEvent) {
			call(e)
		})
	default:
	}
}

func (a *App) Close() {
	a.conn.Close()
	fmt.Println("close")
}
