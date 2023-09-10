package loadbalancer

import (
	"log"
	"net/http"
	"net/url"
)

func InitLoadBalancer(urls []string) error {
	SetupLoadBalancerRoutes()
	go monitorBackendServerHealth()
	return LoadServerPool(urls)
}

func LoadServerPool(urls []string) error {
	for _, v := range urls {
		validURL, err := url.Parse(v)
		if err != nil {
			return err
		}
		ServerPool = append(ServerPool, Server{URL: validURL, Health: true, HealthChan: make(chan bool)})
	}
	log.Println("Server Pool successfully initialized")
	return nil
}

func SetupLoadBalancerRoutes() {
	http.HandleFunc("/", loadBalanceHandler)
}

func NewLoadBalancer(port string) error {
	return http.ListenAndServe(port, nil)
}
