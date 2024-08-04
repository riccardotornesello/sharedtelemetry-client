package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sharedtelemetry/client/common"
	"sharedtelemetry/client/iracing"
	"sharedtelemetry/client/websocket"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

type Data struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func sendDataToClients(hub *websocket.Hub, iRacingConnection *iracing.IRacingConnection, fps int) {
	for {
		start := time.Now()

		drivers, race, telemetry := iRacingConnection.GetData()

		driversData := Data{
			Event: "drivers",
			Data:  drivers,
		}
		driversOutput, err := json.Marshal(driversData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.BroadcastMessage("drivers", driversOutput)
		}

		raceData := Data{
			Event: "race",
			Data:  race,
		}
		raceOutput, err := json.Marshal(raceData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.BroadcastMessage("race", raceOutput)
		}

		telemetryData := Data{
			Event: "telemetry",
			Data:  telemetry,
		}
		telemetryOutput, err := json.Marshal(telemetryData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.BroadcastMessage("telemetry", telemetryOutput)
		}

		elapsed := time.Since(start)
		time.Sleep(time.Second/time.Duration(fps) - elapsed)
	}
}

func main() {
	flag.Parse()

	hub := websocket.NewHub()
	go hub.Run()

	// TODO: emit events
	eventsChannel := make(chan common.Event)

	iRacingConnection := iracing.NewConnection()
	go iRacingConnection.Start(10, 1000, eventsChannel)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	go sendDataToClients(hub, iRacingConnection, 60)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
