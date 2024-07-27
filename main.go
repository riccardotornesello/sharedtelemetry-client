package main

import (
	"encoding/json"
	"example/sharedtelemetry/iracing"
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

type Data struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func sendDataToClients(hub *Hub, iRacingConnection *iracing.IRacingConnection) {
	for {
		drivers, event, telemetry := iRacingConnection.GetData()

		driversData := Data{
			Event: "drivers",
			Data:  drivers,
		}
		driversOutput, err := json.Marshal(driversData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.broadcast <- driversOutput
		}

		eventData := Data{
			Event: "event",
			Data:  event,
		}
		eventOutput, err := json.Marshal(eventData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.broadcast <- eventOutput
		}

		telemetryData := Data{
			Event: "telemetry",
			Data:  telemetry,
		}
		telemetryOutput, err := json.Marshal(telemetryData)
		if err != nil {
			log.Println("Error marshalling data: ", err)
		} else {
			hub.broadcast <- telemetryOutput
		}

		time.Sleep(20 * time.Millisecond)
	}
}

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()

	iRacingConnection := iracing.NewConnection()
	go iRacingConnection.Start(10, 1000)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go sendDataToClients(hub, iRacingConnection)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
