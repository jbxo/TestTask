package shared

import (
	"net/http"
)

type Adapter func(handler http.Handler) http.Handler

// adaptersHandlerExecutor wraps handler with Adapters from the slice in reverse order.
func adaptersHandlerExecutor(h http.Handler, as []Adapter) http.Handler {
	for i := len(as) - 1; i >= 0; i-- {
		h = as[i](h)
	}

	return h
}
