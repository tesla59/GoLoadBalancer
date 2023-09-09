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
	WorkerID            string `gorm:"primaryKey"`
	Timestamp           time.Time
	SuccessfulRequests  int
	FailedRequests      int
	TotalRequests       int
	AverageResponseTime float64
}
