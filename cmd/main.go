package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	algorithms "goLB/algo"
	"goLB/healthcheck"
	backend "goLB/utility"
)

var (
	backends []*backend.Backend // Slice of backend servers
	mutex    sync.Mutex
)

var (
	algos           = backend.WeightedRoundRobbin
	nextServerIndex = 0
	CurrentWeight   = 0
)

func main() {

	data, err1 := backend.ReadConfig("config.json")
	if err1 != nil {
		fmt.Println("Error reading config file", err1)
	}

	urlString := data.Proxy

	parsedURL, err2 := url.Parse(urlString)
	if err2 != nil {
		fmt.Println("Error parsing URL:", err2)
		return
	}

	_, port, err3 := net.SplitHostPort(parsedURL.Host)
	if err3 != nil {
		fmt.Println("Error parsing port:", err3)
		return
	}

	strPort := fmt.Sprintf(":%s", port)

	backends = make([]*backend.Backend, len(data.Servers))

	for i, s := range data.Servers {
		backends[i] = &backend.Backend{
			URL:         s.URL,
			Weight:      s.Weight,
			Healthy:     true,
			Connections: 0,
		}
		fmt.Println("Backend server", backends[i])
	}

	// Starting HTTP server
	http.HandleFunc("/", handler)

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GoScale is running successfully OK")
	})

	fmt.Printf("Load balancer started on %s", port)
	log.Fatal(http.ListenAndServe(strPort, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	backend := selectHealthyBackend()

	if backend == nil {
		fmt.Println("No healthy backend servers available")
		http.Error(w, "No healthy backend servers available", http.StatusServiceUnavailable)
		return
	}

	fmt.Println("Selected backend server:", backend.URL)
	url, _ := url.Parse(backend.URL)

	// Increment the number of connections for the selected backend server
	backend.IncrementConnections()

	// Forward request to selected backend server
	startTime := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConns:        100, // connection alive until 100s
		MaxIdleConnsPerHost: 100,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	proxy.ServeHTTP(w, r)

	backend.ResponseTime = time.Since(startTime)
	fmt.Println("Response time", backend.ResponseTime)

	// Decrement the number of connections for the selected backend server
	defer backend.DecrementConnections()
}

func selectHealthyBackend() *backend.Backend {
	// Filter out unhealthy backends
	healthyBackends := make([]*backend.Backend, 0)
	for i := range backends {
		if backends[i].Healthy {
			healthyBackends = append(healthyBackends, backends[i])
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	switch algos {
	case backend.RoundRobbin:
		nextServerIndex = algorithms.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.LeastConnections:
		nextServerIndex = algorithms.LeastConnections(healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.WeightedRoundRobbin:
		// fmt.Println("CurrentWeight is the", CurrentWeight)
		nextServerIndex, CurrentWeight = algorithms.WeightedRoundRobbin(CurrentWeight, healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.LeastTime:
		nextServerIndex = algorithms.LeastTime(healthyBackends)
		return healthyBackends[nextServerIndex]

	default:
		nextServerIndex = algorithms.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]

	}
}

func healthCheck() {
	for {
		time.Sleep(10 * time.Second) // Check health every 10 seconds

		// Perform health check for each backend server
		for i := range backends {
			if !healthcheck.CheckHealth(backends[i]) {
				backends[i].Healthy = false
				fmt.Printf("Backend %s is unhealthy %d \n", backends[i].URL, backends[i].Connections)
			} else {
				backends[i].Healthy = true
				fmt.Printf("Backend %s is healthy %d \n", backends[i].URL, backends[i].Connections)
			}
		}
	}
}

func init() {
	// Start health check in the background
	go healthCheck()
}
