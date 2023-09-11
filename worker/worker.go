package main

import (
	"net/http"
	"os"
)

var HostName string
var config Config

func init() {
	HostName, _ = os.Hostname()
}

func main() {
	InitDB()
	readYAML()

	http.HandleFunc("/api/v1/hello", HelloResponse)
	http.HandleFunc("/worker/stats", GetWorkerStats)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
