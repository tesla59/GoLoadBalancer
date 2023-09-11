package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type StatsResponse struct {
	SuccessfulRequests  []map[string]int
	FailedRequests      []map[string]int
	TotalRequests       []map[string]int
	AverageResponseTime []map[string]float64
}

type Config struct {
	Worker   int     `yaml:"worker"`
	Pool     int     `yaml:"pool"`
	StatsDir string  `yaml:"stats-dir"`
	AvgDelay float64 `yaml:"avg-delay"`
	Failure  int     `yaml:"failure"`
}

func HelloResponse(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	Worker := WorkerStats{
		WorkerID: HostName,
	}
	DB.First(&Worker)

	message := map[string]string{"message": "hello from container: " + HostName}
	messageJSON, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")

	isASuccessfulRequest := rand.Intn(100) > config.Failure

	if isASuccessfulRequest {
		Worker.SuccessfulRequests++
		// Average Delay
		Sleep(int(config.AvgDelay))

		if _, err := w.Write(messageJSON); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		Worker.FailedRequests++
		http.Error(w, "Randomized Server Error", http.StatusInternalServerError)
	}
	Worker.AverageResponseTime = float64(time.Since(start).Milliseconds())
	Worker.TotalRequests = Worker.SuccessfulRequests + Worker.FailedRequests
	DB.Save(&Worker)
}

func GetWorkerStats(w http.ResponseWriter, r *http.Request) {
	var Stats []WorkerStats
	var statsResponse StatsResponse

	DB.Find(&Stats)
	totalResponseTime, totalFailedReq, totalSuccessfulReq, totalTotalReq := float64(0), 0, 0, 0

	// Setting total fields first
	for i := range Stats {
		totalResponseTime += Stats[i].AverageResponseTime
		totalFailedReq += Stats[i].FailedRequests
		totalSuccessfulReq += Stats[i].SuccessfulRequests
		totalTotalReq += Stats[i].TotalRequests
	}

	statsResponse.AverageResponseTime = append(statsResponse.AverageResponseTime, map[string]float64{"total": totalResponseTime})
	statsResponse.FailedRequests = append(statsResponse.FailedRequests, map[string]int{"total": totalFailedReq})
	statsResponse.SuccessfulRequests = append(statsResponse.SuccessfulRequests, map[string]int{"total": totalSuccessfulReq})
	statsResponse.TotalRequests = append(statsResponse.TotalRequests, map[string]int{"total": totalTotalReq})

	// Setting worker fields later
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
