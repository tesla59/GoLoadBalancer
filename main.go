package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func init() {
	BuildImage("./worker", "worker-image")

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
}

var config Config

func main() {
}
