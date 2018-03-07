package shared

import (
	"log"
	"net/http"
)

func HandleErrors() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					error, ok := rec.(Error)
					// If the panic argument object is a type of shared.Error
					// Then handle it
					// Else pass it forward
					if ok {
						handleError(w, error)
					} else {
						panic(rec)
					}
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

func handleError(w http.ResponseWriter, error Error) {
	log.Printf("[ERROR] An error occured during query execution. InternalMessage: %v."+
		" Details: %v. Client will get HTTP-Code: %v. Message: %v.\n",
		error.InternalMessage, error.Details, error.Code, error.UserMessage)

	// Just pass it to http error mech.
	http.Error(w, error.UserMessage, error.Code)
}
