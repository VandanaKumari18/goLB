package algorithms

import (
	backend "goLB/utility"
)

func RoundRobbin(index int, backends []*backend.Backend) int {
	// fmt.Println(backends[0].ResponseTime)
	return (index + 1) % len(backends)

}
