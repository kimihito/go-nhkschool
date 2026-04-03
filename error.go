package nhkschool

import "fmt"

// APIError is returned when the API responds with a non-2xx status code.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Body is the raw response body.
	Body string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("nhkschool: HTTP %d: %s", e.StatusCode, e.Body)
}
