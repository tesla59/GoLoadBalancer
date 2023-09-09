package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func init() {
	BuildImage("./worker", "worker-image")
	PWD, _ = os.Getwd()

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
}

var config Config
var PWD string

func main() {
}
