package main

import (
	"github.com/zing-dev/dts-sdk/sdk/client"
	"log"
	"time"
)

func main() {
	app := client.New(client.Option{Ip: "192.168.0.215", Port: 17083})
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
	app.Run()
	time.Sleep(time.Minute * 5)
}
