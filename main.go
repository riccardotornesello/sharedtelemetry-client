package main

import (
	"flag"
	"log"
	"net/http"
	"sharedtelemetry/client/common"
	"sharedtelemetry/client/iracing"
	"sharedtelemetry/client/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

type Data struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func main() {
	flag.Parse()

	eventsChannel := make(chan common.Event)

	hub := websocket.NewHub()
	go hub.Run(eventsChannel)

	iRacingConnection := iracing.NewConnection()
	iRacingConnection.Start(eventsChannel, 60, 10)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
