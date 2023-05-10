package api_error

type ApiError struct {
	ErrorMessage string
	ErrorCode    string
	StatusCode   int
}

func (f *ApiError) Error() string {
	return f.ErrorMessage
}
