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

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func sendDataToClients(hub *Hub, iRacingConnection *iracing.IRacingConnection) {
	for {
		drivers, _ := iRacingConnection.GetData()
		driversData, err := json.Marshal(drivers)
		if err != nil {
			log.Println("Error marshalling drivers data: ", err)
		} else {
			hub.broadcast <- driversData
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()

	iRacingConnection := iracing.NewConnection()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go sendDataToClients(hub, iRacingConnection)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
