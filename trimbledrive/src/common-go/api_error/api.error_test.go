package api_error

import (
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	filespaceError := &ApiError{
		ErrorMessage: "Please provide a valid value",
		ErrorCode:    "InvalidPayload",
		StatusCode:   http.StatusBadRequest,
	}
	got := filespaceError.Error()
	if got != "Please provide a valid value" {
		t.Errorf("Unexpected result `%s`", got)
	}
}
