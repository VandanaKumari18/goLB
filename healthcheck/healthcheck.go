package healthcheck

import (
	backend "goLB/utility"
	"net/http"
)

func CheckHealth(backend *backend.Backend) bool {
	// Send a health check request to the backend server
	resp, err := http.Get(backend.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
