package nhkschool

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Body: "Not Found"}
	want := "nhkschool: HTTP 404: Not Found"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_ErrorsAs(t *testing.T) {
	orig := &APIError{StatusCode: 429, Body: "Rate limited"}
	var wrapped error = orig

	var apiErr *APIError
	if !errors.As(wrapped, &apiErr) {
		t.Fatal("errors.As failed to match *APIError")
	}
	if apiErr.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
}
