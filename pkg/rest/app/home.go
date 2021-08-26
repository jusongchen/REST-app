package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// home returns a simple HTTP handler function which writes a response.
func home(a Instance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		info := struct {
			UpTime      string `json:"up_time,omitempty"`
			StartupTime time.Time
			Host        string
			Port        uint
			SwaggerDir  string
		}{

			time.Now().Sub(a.StartupTime).String(),
			a.StartupTime,
			a.Config.Host,
			a.Config.Port,
			a.Config.SwaggerDir,
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			log.Printf("Could not encode info data: %v", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain")

		io.WriteString(w, string(data)+"\n"+a.Config.About)
	}
}
