package main

import "time"

type Config struct {
	Worker           int     `yaml:"worker"`
	Pool             int     `yaml:"pool"`
	StatsDir         string  `yaml:"stats-dir"`
	AvgDelay         float64 `yaml:"avg-delay"`
	Failure          int     `yaml:"failure"`
	LoadBalancerPort int     `yaml:"load-balancer-port"`
	WorkerPort       int     `yaml:"worker-port"`
}

type WorkerStats struct {
	WorkerID            string `gorm:"primaryKey"`
	Timestamp           time.Time
	SuccessfulRequests  int
	FailedRequests      int
	TotalRequests       int
	AverageResponseTime float64
}
