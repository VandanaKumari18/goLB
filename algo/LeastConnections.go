package algorithms

import (
	backend "goLB/utility"
)

func LeastConnections(backends []*backend.Backend) int {
	var min int
	min = int(backends[0].Connections)
	index := 0
	for i := 0; i < len(backends); i++ {
		if int(backends[i].Connections) < min {
			min = int(backends[i].Connections)
		}
	}
	return index

}
