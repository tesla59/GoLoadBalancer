package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

func init() {
	BuildImage(WorkerSrcPath, DefaultProtocol)
	PWD, _ = os.Getwd()

	yamlFile, err := os.ReadFile(DefaultConfigFilePath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}
}

func main() {
}
