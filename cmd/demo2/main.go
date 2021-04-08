package main

import (
	"github.com/zing-dev/dts-sdk/sdk/client"
	"github.com/zing-dev/dts-sdk/sdk/msg/models"
	"log"
	"time"
)

func main() {
	app := client.New(client.Option{Ip: "192.168.0.215", Port: 17083})
	/*app.SetZoneTempNotify(func(notify *models.ZonesTemp) {
		z := new(models.ZonesTemp)
		z.CreatedAt = notify.CreatedAt
		z.DeviceId = notify.DeviceId
		z.Zones = notify.Zones[:10]
		log.Println("ZoneTempNotify", len(notify.Zones))
		data, _ := json.Marshal(z)
		fmt.Println(string(data))
	})*/
	app.SetTempSignalNotify(func(notify *models.TempSignalNotify) {
		log.Println("TempSignalNotify", len(notify.Signal))
	})
	app.Subscribe(client.TopicTemp, func() {})
	app.Run()
	//app.Run()
	time.Sleep(time.Minute * 5)
}
