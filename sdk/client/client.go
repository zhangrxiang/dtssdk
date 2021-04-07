package client

import (
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

type App struct {
	Option
	EnableAlarm bool
	EnableSign  bool
	EnableTemp  bool
	EnableEvent bool

	conn  *tao.ClientConn
	group *sync.WaitGroup
}

func New(o Option) *App {
	return &App{
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
			err := a.conn.Write(&request.PingRequest{})
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
	fmt.Println("success")
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

func (a *App) SetTempSignalNotify(f func(*models.TempSignalNotify)) {
	a.group.Add(1)
	defer func() {
		a.group.Done()
		fmt.Println("sign done")
	}()

	s := response.SignResponse{
		Request: new(models.TempSignalNotify),
		Device:  new(models.SetDeviceRequest),
		Value:   make(chan *models.TempSignalNotify, 10),
	}
	tao.Register(s.MessageNumber(), s.Unmarshaler, s.Handle)
	go func(f func(*models.TempSignalNotify)) {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case v := <-s.Value:
				t.Reset(time.Minute)
				f(v)
			case <-t.C:
				a.conn.Close()
				return
			}
		}
	}(f)
}

func (a *App) SetZoneTempNotify(f func(*models.ZoneTempNotify)) {
	a.group.Add(1)
	defer func() {
		a.group.Done()
		fmt.Println("temp done")
	}()
	s := response.TempRequest{
		Request: new(models.ZoneTempNotify),
		Device:  new(models.SetDeviceRequest),
		Value:   make(chan *models.ZoneTempNotify, 10),
	}
	tao.Register(s.MessageNumber(), s.Unmarshaler, s.Handle)
	go func(f func(*models.ZoneTempNotify)) {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case v := <-s.Value:
				t.Reset(time.Minute)
				f(v)
			case <-t.C:
				a.conn.Close()
			}
		}
	}(f)
}

func (a *App) Close() {
	a.conn.Close()
	fmt.Println("close")
}
