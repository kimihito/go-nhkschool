package nhkschool

import "fmt"

// APIError is returned when the API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("nhkschool: HTTP %d: %s", e.StatusCode, e.Body)
}
