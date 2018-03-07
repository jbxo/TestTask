package shared

import (
	"log"
	"net/http"
	"time"
)

func LogQueries() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Printf("[INFO] Request to %v finished in %v", r.URL.String(), time.Since(start))
		})
	}
}
