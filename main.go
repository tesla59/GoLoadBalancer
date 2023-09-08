package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func init() {
	BuildImage("./worker", "worker-image")
}

var config Config

func main() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
}
