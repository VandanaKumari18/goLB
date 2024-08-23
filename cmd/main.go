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

	backend "goLB/constants"
)

var backends []*backend.Backend // Slice of backend servers
var mutex sync.Mutex

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

	backend := backends[0]

	if backend == nil {
		fmt.Println("No healthy backend servers available")
		http.Error(w, "No healthy backend servers available", http.StatusServiceUnavailable)
		return
	}

	fmt.Println("Selected backend server:", backend.URL)
	url, _ := url.Parse(backend.URL)

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

}
