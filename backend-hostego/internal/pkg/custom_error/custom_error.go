package customerror

type CustomError struct {
	Message    string
	Code       string
	ErrorValue error
	StatusCode int
}

func NewCustomError(message, code string, errorValue error, statusCode int) *CustomError {
	return &CustomError{
		Message:    message,
		Code:       code,
		ErrorValue: errorValue,
		StatusCode: statusCode,
	}
}

func NewNilCustomError() *CustomError {
	return &CustomError{
		ErrorValue: nil,
	}
}
