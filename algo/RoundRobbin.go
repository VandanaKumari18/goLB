package algorithms

import "goLB/constants"

func RoundRobbin(index int, backends []*constants.Backend) int {
	// fmt.Println(backends[0].ResponseTime)
	return (index + 1) % len(backends)

}
