package constants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Servers []ServerConfig `json:"servers"`
	Proxy   string         `json:"proxy"`
}

type ServerConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

type Backend struct {
	URL           string `json:"url"`
	Healthy       bool
	Connections   int
	ResponseTime  time.Duration
	Weight        int `json:"weight"`
	WeightedScore float64
}

var (
	RoundRobbin         = "RoundRobbin"
	LeastConnections    = "LeastConnections"
	LeastTime           = "LeastTime"
	WeightedRoundRobbin = "WeightedRoundRobbin"
)

// util methods

func (b *Backend) IncrementConnections() {
	b.Connections++
	fmt.Println("Connections increment", b.Connections)
}
func (b *Backend) DecrementConnections() {
	fmt.Println("Connections decrement", b.Connections)
	b.Connections--
}

func ReadConfig(filename string) (*Config, error) {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close() // To parse json later on, we'll defer its closing

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	var backends []Backend // Populate Backend struct array

	for _, server := range config.Servers {
		backend := Backend{
			URL:     server.URL,
			Weight:  server.Weight,
			Healthy: true, // healthy server initialization
		}
		backends = append(backends, backend)
	}

	return &config, nil
}
