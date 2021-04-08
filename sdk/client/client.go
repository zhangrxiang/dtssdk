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
	"time"
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
	conn  *tao.ClientConn
	group *sync.WaitGroup
}

func New(o Option) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
		Option: o,
		group:  new(sync.WaitGroup),
	}
}

func (a *App) valid() {
	if a.Option.Ip == "" {
		panic("ip")
	}

	if a.Option.Port < 1024 {
		panic("port")
	}
}

func (a *App) Ping() {
	go func() {
		a.group.Add(1)
		defer func() {
			fmt.Println("ping done")
			a.group.Done()
		}()
		for {
			err := a.conn.Write(&request.PingReq{})
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

func (a *App) Run() {
	a.valid()
	fmt.Println("start", fmt.Sprintf("%s:%d", a.Ip, a.Port))
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.Ip, a.Port))
	if err != nil {
		log.Fatalln(err)
	}
	onConnect := tao.OnConnectOption(func(conn tao.WriteCloser) bool {
		a.Ping()
		_ = conn.Write(&request.DeviceRequest{Request: &models.SetDeviceRequest{
			ZoneTempNotifyEnable:    true,
			ZoneAlarmNotifyEnable:   false,
			FiberStatusNotifyEnable: false,
			TempSignalNotifyEnable:  true,
		}})
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

func (a *App) Publish(topic TopicType) {
	switch topic {

	}
}
func (a *App) Subscribe(topic TopicType, call func(result interface{})) {
	switch topic {
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
