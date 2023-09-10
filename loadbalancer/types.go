package loadbalancer

import (
	"net/url"
	"sync"
)

var ServerPool []Server
var currentIndex int
var mutex sync.Mutex

type Server struct {
    URL        *url.URL
    Health     bool
    Mutex      sync.Mutex
    HealthChan chan bool
}
