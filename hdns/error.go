package hdns

import "fmt"

// ErrorCode represents an error code returned from the API.
type ErrorCode int

// Error codes returned from the API.
const (
	ErrorCodeNotFound          ErrorCode = 404
	ErrorCodeRateLimitExceeded ErrorCode = 429
)

// Error is an error returned from the API.
type Error struct {
	Code    ErrorCode
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// IsError returns whether err is an API error with the given error code.
func IsError(err error, code ErrorCode) bool {
	apiErr, ok := err.(Error)
	return ok && apiErr.Code == code
}
