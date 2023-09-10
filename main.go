package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gocomet/loadbalancer"
)

func init() {
	BuildImage(WorkerSrcPath, DefaultWorkerImageTag)
	PWD, _ = os.Getwd()

	yamlFile, err := os.ReadFile(DefaultConfigFilePath)
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
	RunImage(DefaultWorkerImageTag, 10, 8080)
	loadbalancer.InitLoadBalancer([]string{"http://localhost:8080", "http://localhost:8081"})
	go loadbalancer.NewLoadBalancer(":8000")

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
		os.Exit(0)
	}()
}
