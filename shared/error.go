package shared

import "net/http"

type DefError struct {
	UserMessage string
	Code        int
}

type Error struct {
	DefError
	InternalMessage string
	Details         error
}

var (
	ErrNotFound         = NewDefError("Not found", http.StatusNotFound)
	ErrMethodNotAllowed = NewDefError("Method not allowed", http.StatusMethodNotAllowed)
	ErrInternalError    = NewDefError("Internal server error", http.StatusInternalServerError)
	ErrBadRequest       = NewDefError("Bad request", http.StatusBadRequest)
)

func NewDefError(uMessage string, code int) DefError {
	return DefError{uMessage, code}
}

func ThrowError(dError DefError, intMessage string, err error) {
	panic(Error{dError, "An error occurred while " + intMessage, err})
}
