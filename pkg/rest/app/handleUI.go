package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jusongchen/REST-app/pkg/logging"
	"github.com/rs/cors"
)

// StaticHTMLHandler returns a HandlerFunc which serves static html files
// Enhancement:
//   1) Take allowCORS as a parameter
//	 2) if the directory to store the static html files does not exists, return 500 with explicit error message
//		instread of return StatusServiceUnavailable
//
func StaticHTMLHandler(urlPath string, staticFilePath string, allowCORS bool) http.HandlerFunc {

	logger := logging.FromContext(context.Background()).Named("StaticHTMLHandler")

	exist, err := exists(staticFilePath)
	if exist {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("Request:from %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
			h := http.StripPrefix(urlPath, http.FileServer(http.Dir(staticFilePath)))
			if allowCORS {
				cors.Default().Handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		}
	}

	if err == nil && !exist {
		cwd, _ := os.Getwd()
		err = fmt.Errorf("serve Static HTML: (Current Working Dir %s) directory does not exist:%s", cwd, staticFilePath)
	}
	logger.Errorw("error init static html handler", "error", err)
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
