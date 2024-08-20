package constants

import "time"

type Backend struct {
	URL           string `json:"url"`
	Healthy       bool
	Connections   int
	ResponseTime  time.Duration
	Weight        int `json:"weight"`
	WeightedScore float64
}
