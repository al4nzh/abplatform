package main

import (
	"bytes"
	"encoding/json"
	//"log"
	"math/rand"
	"net/http"
	"time"
)

type Event struct {
	ExperimentID int    `json:"experiment_id"`
	UserID       string `json:"user_id"`
	EventName    string `json:"event_name"`
}

func main() {
	for i := 0; i < 1000; i++ {
		user := "user" + string(rune(rand.Intn(1000)))

		// assign
		http.Get("http://localhost:8080/assign?experiment_id=1&user_id=" + user)

		// impression
		sendEvent(user, "impression")

		// random conversion
		if rand.Float64() < 0.3 {
			sendEvent(user, "click")
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func sendEvent(user, eventName string) {
	event := Event{
		ExperimentID: 1,
		UserID:       user,
		EventName:    eventName,
	}

	data, _ := json.Marshal(event)

	http.Post("http://localhost:8080/events", "application/json", bytes.NewBuffer(data))
}