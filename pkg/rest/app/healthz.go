package app

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// healthz is a liveness probe.
func healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request:from %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("service is alive.\n"))
	}
}
