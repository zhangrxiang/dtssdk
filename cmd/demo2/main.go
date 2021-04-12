package main

import (
	"github.com/zing-dev/dts-sdk/sdk/client"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"github.com/zing-dev/dts-sdk/sdk/msg/request"
	"log"
	"time"
)

func main() {
	app := client.New(client.Option{Ip: "192.168.0.215", Port: 17083, ChannelNum: 1})
	app.Subscribe(client.TopicTemp, func(result interface{}) {
		log.Println("temp")
	})
	app.Subscribe(client.TopicSignal, func(result interface{}) {
		log.Println("signal")
	})
	app.Subscribe(client.TopicAlarm, func(result interface{}) {
		log.Println("alarm")
	})
	app.Subscribe(client.TopicEvent, func(result interface{}) {
		log.Println("event")
	})
	app.Subscribe(client.TopicZones, func(result interface{}) {
		log.Println("zones", len(result.(*models.Zones).Zones))
	})
	app.Publish(request.NewZones(4))
	app.Run()
	time.Sleep(time.Minute * 5)
}
