package model

// ErrorResponse defines model for Error Response Body.
type ErrorResponse struct {
	// Code The error code. A non-empty language-independent string.
	Code string `json:"code"`
	// Message A human-readable string describing the errors.
	Message string `json:"message"`
}

func New(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
