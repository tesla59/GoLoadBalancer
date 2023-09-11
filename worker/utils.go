package main

import (
	"math/rand"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

func readYAML() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
}

func Sleep(averageResponseTimeMs int) {
	minSleepMs := averageResponseTimeMs / 2
	maxSleepMs := averageResponseTimeMs * 3 / 2
	sleepMs := rand.Intn(maxSleepMs-minSleepMs+1) + minSleepMs
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}
