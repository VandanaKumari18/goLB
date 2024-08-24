package utility

import (
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
