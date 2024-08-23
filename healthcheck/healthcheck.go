package healthcheck

import (
	"goLB/constants"
	"net/http"
)

func CheckHealth(backend *constants.Backend) bool {
	// Send a health check request to the backend server
	resp, err := http.Get(backend.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
