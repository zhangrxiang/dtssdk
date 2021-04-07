package main

import (
	"github.com/zing-dev/dts-sdk/sdk/client"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"log"
	"time"
)

func main() {
	app := client.New(client.Option{Ip: "192.168.0.215", Port: 17083})
	app.SetZoneTempNotify(func(notify *models.ZoneTempNotify) {
		log.Println("ZoneTempNotify", len(notify.Zones))
	})
	app.SetTempSignalNotify(func(notify *models.TempSignalNotify) {
		log.Println("TempSignalNotify", len(notify.Signal))
	})
	go app.Run()
	time.Sleep(time.Second * 10)
	app.Close()
	time.Sleep(time.Second * 10)
	//app.Run()
	//time.Sleep(time.Minute*5)
}
