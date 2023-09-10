package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
    backendServer1, _ := url.Parse("http://localhost:8080")
    backendServer2, _ := url.Parse("http://localhost:8081")

    ServerPool = append(ServerPool, Server{URL: backendServer1, Health: true, HealthChan: make(chan bool)})
    ServerPool = append(ServerPool, Server{URL: backendServer2, Health: true, HealthChan: make(chan bool)})

	http.HandleFunc("/", loadBalanceHandler)

	// Start the load balancer server.
    if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func loadBalanceHandler(w http.ResponseWriter, r *http.Request) {
	backendServer := selectBackendServer()

	// Proxy the request to the selected backend server.
	proxy := createReverseProxy(backendServer.URL)
	proxy.ServeHTTP(w, r)
}

func selectBackendServer() *Server {
	mutex.Lock()
	defer mutex.Unlock()

	// Find a healthy backend server using round-robin.
	for i := 0; i < len(ServerPool); i++ {
		currentIndex = (currentIndex + 1) % len(ServerPool)
		server := &ServerPool[currentIndex]
		if server.Health {
			return server
		}
	}

	// No healthy server found, use the first one as a fallback.
	return &ServerPool[0]
}

func createReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(target)
}
