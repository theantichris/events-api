package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Custom HTTP errors.
var (
	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong.",
	}

	ErrUnprocessableEntity = &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "Unprocessable entity.",
	}

	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Error invalid argument.",
	}

	ErrEventNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Event not found.",
	}

	ErrObjectIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Request object should be provided.",
	}

	ErrValidEventIDIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "A valid event ID is required.",
	}

	ErrEventTimingIsRequired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Event start time and end time are required.",
	}

	ErrInvalidLimit = &Error{
		Code:    http.StatusBadRequest,
		Message: "Limit should be an integral value.",
	}

	ErrInvalidTimeFormat = &Error{
		Code:    http.StatusBadRequest,
		Message: "Time should be passed in RFC3339 Format: " + time.RFC3339,
	}
)

// Error holds information on any errors that occur.
type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) String() string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

// JSON serializes an error into JSON.
func (err *Error) JSON() []byte {
	if err == nil {
		return []byte("{}")
	}

	res, _ := json.Marshal(err)

	return res
}

// StatusCode returns the HTTP status code for the error.
func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}

	return err.Code
}
