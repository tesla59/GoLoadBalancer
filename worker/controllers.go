package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type StatsResponse struct {
	SuccessfulRequests  []map[string]int
	FailedRequests      []map[string]int
	TotalRequests       []map[string]int
	AverageResponseTime []map[string]float64
}

func HelloResponse(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{"message": "hello from container: " + HostName}
	messageJSON, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")

	Worker := WorkerStats{
		WorkerID: HostName,
	}

	DB.First(&Worker)
	Worker.SuccessfulRequests++
	Worker.TotalRequests = Worker.SuccessfulRequests + Worker.FailedRequests
	DB.Save(&Worker)

	// Add a delay of 500 ms to replicate some labor intensive task
	time.Sleep(time.Millisecond * 500)

	if _, err := w.Write(messageJSON); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetWorkerStats(w http.ResponseWriter, r *http.Request) {
	var Stats []WorkerStats
	var statsResponse StatsResponse

	DB.Find(&Stats)

	for i := range Stats {
		statsResponse.AverageResponseTime = append(statsResponse.AverageResponseTime, map[string]float64{
			Stats[i].WorkerID: Stats[i].AverageResponseTime,
		})
		statsResponse.FailedRequests = append(statsResponse.FailedRequests, map[string]int{
			Stats[i].WorkerID: Stats[i].FailedRequests,
		})
		statsResponse.SuccessfulRequests = append(statsResponse.SuccessfulRequests, map[string]int{
			Stats[i].WorkerID: Stats[i].SuccessfulRequests,
		})
		statsResponse.TotalRequests = append(statsResponse.TotalRequests, map[string]int{
			Stats[i].WorkerID: Stats[i].TotalRequests,
		})
	}
	statsResponseJSON, err := json.Marshal(statsResponse)
	if err != nil {
		http.Error(w, "err: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(statsResponseJSON); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
