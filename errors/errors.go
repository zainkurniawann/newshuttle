package errors

type CustomError struct {
    Message    string
    StatusCode int
}

func New(message string, statusCode int) *CustomError {
    return &CustomError{
        Message:    message,
        StatusCode: statusCode,
    }
}

func (e *CustomError) Error() string {
    return e.Message
}