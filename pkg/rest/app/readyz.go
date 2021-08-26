package app

import (
	"net/http"
	"sync/atomic"
)

// readyz is a readiness probe.
func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(http.StatusText(http.StatusServiceUnavailable)))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
