package main

import (
	"os"

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
	loadbalancer.InitLoadBalancer([]string{"http://localhost:8080", "http://localhost:8081"})
	loadbalancer.NewLoadBalancer(":8000")
}
