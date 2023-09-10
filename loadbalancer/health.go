package main

import (
	"net/http"
	"time"
)

func monitorBackendServerHealth() {
    for {
        time.Sleep(5 * time.Second)
        for i := 0; i < len(ServerPool); i++ {
            server := &ServerPool[i]
            server.Mutex.Lock()
            _, err := http.Head(server.URL.String())
            if err != nil {
                server.Health = false
            } else {
                server.Health = true
            }
            server.Mutex.Unlock()
        }
    }
}
