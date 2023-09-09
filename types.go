package main

import "time"

type Config struct {
	Worker    int     `yaml:"worker"`
	Pool      int     `yaml:"pool"`
	StatsDir  string  `yaml:"stats-dir"`
	AvgDelay  float64 `yaml:"avg-delay"`
	Failure   int     `yaml:"failure"`
}


type WorkerStats struct {
	WorkerID           int       `db:"worker_id"`
	Timestamp          time.Time `db:"timestamp"`
	SuccessfulRequests int       `db:"successful_requests"`
	FailedRequests     int       `db:"failed_requests"`
	TotalRequests      int       `db:"total_requests"`
	AverageResponseTime float64   `db:"average_response_time"`
}
