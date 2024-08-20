package algorithms

import (
	"goLB/constants"
)

func LeastConnections(backends []*constants.Backend) int {
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
