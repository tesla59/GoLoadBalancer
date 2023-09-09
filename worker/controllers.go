package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func HelloResponse(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{"message": "hello from " + HostName}

	messageJSON, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")

	// Add a delay of 500 ms to replicate some labor intensive task
	time.Sleep(time.Millisecond * 500)

	if _, err := w.Write(messageJSON); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
