package shared

import "net/http"

func DefineHTTPMethod(method string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				http.Error(w, "That method isn't supported.", http.StatusBadRequest)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
