package loadbalancer

import (
	"errors"
	"net/http"
	"net/http/httputil"
)

func loadBalanceHandler(w http.ResponseWriter, r *http.Request) {
	backendServer, err := selectBackendServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Proxy the request to the selected backend server.
	proxy := httputil.NewSingleHostReverseProxy(backendServer.URL)
	proxy.ServeHTTP(w, r)
}

func selectBackendServer() (*Server, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// Find a healthy backend server using round-robin.
	for i := 0; i < len(ServerPool); i++ {
		currentIndex = (currentIndex + 1) % len(ServerPool)
		server := &ServerPool[currentIndex]
		if server.Health {
			return server, nil
		}
	}

	// No healthy server found, return error
	return nil, errors.New("No healthy server found")
}
