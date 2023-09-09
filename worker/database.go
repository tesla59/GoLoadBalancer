package main

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type WorkerStats struct {
	WorkerID            string `gorm:"primaryKey"`
	Timestamp           time.Time
	SuccessfulRequests  int
	FailedRequests      int
	TotalRequests       int
	AverageResponseTime float64
}

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("worker-stats.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	Worker := WorkerStats{
		WorkerID:            HostName,
		Timestamp:           time.Now(),
		SuccessfulRequests:  0,
		FailedRequests:      0,
		TotalRequests:       0,
		AverageResponseTime: 0,
	}

	DB.Create(&Worker)
}
