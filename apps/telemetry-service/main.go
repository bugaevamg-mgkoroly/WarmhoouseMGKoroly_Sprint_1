package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/streadway/amqp"
)

type Telemetry struct {
	DeviceID  string  `json:"device_id"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp,omitempty"`
}

var conn *amqp.Connection

func initRMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil { log.Fatalf("RMQ conn: %v", err) }
}

func publish(data Telemetry) {
	ch, err := conn.Channel()
	if err != nil { return }
	defer ch.Close()

	body, _ := json.Marshal(data)
	err = ch.Publish("", "telemetry", false, false, amqp.Publishing{
		ContentType: "application/json", Body: body })
	if err != nil { log.Printf("Pub err: %v", err) }
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data Telemetry
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.DeviceID == "" {
		http.Error(w, "Invalid data", 400)
		return
	}
	publish(data)
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
}

func main() {
	initRMQ()
	http.HandleFunc("/telemetry", handler)
	log.Println("Telemetry: :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}