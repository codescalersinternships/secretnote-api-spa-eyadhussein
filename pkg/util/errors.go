package util

import (
	"errors"
	"fmt"
)

var (
	ErrInternalServer = errors.New("internal server error") // internal server error
	ErrBadRequest     = errors.New("bad request")           // a bad request error
	ErrUnauthorized   = errors.New("unauthorized")          // unauthorized user
	ErrNotFound       = errors.New("not found")             // resource not found

	ErrMaxViewsLessThanOne = errors.New("max views has to be at least 1")          // max views is less than 1
	ErrExpiresAtBeforeNow  = errors.New("expiration date has to be in the future") // expiration date is before the current date
)

// ResponseError struct holds the error message and status code
// @Description General error response
// @Response
type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// NewResponseError creates a new response error
func NewResponseError(err error, status int) *ResponseError {
	return &ResponseError{
		Message: err.Error(),
		Status:  status,
	}
}

// Error returns the error message
func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s with status code %d", e.Message, e.Status)
}
