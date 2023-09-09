package main

import (
	"net/http"
	"os"
)

var HostName string

func init() {
	HostName, _ = os.Hostname()
}

func main() {
	http.HandleFunc("/", HelloResponse)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
