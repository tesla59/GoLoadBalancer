package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	BuildImage(WorkerSrcPath, DefaultWorkerImageTag)
	PWD, _ = os.Getwd()

	yamlFile, err := os.ReadFile(DefaultConfigFileName)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(DatabaseFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&WorkerStats{})
}

func main() {
	// Spawn Worker Images
	RunImage(DefaultWorkerImageTag, config.Worker, 8080)
	// Spawn Load Balancers
	RunLoadBalancers(config.Pool, 8000)

	go gracefulShutdown()
	forever := make(chan int)
	<-forever
}

func gracefulShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		log.Println("Shutting down gracefully.")
		StopContainers(ContainerIDs)
		RemoveContainers(ContainerIDs)
		os.Exit(0)
	}()
}
