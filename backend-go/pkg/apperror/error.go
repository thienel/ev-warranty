package apperror

import "fmt"

type AppError struct {
	HttpCode  int
	ErrorCode string
	Message   string
	Err       error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.ErrorCode, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(httpCode int, errorCode, message string) *AppError {
	return &AppError{
		HttpCode:  httpCode,
		ErrorCode: errorCode,
		Message:   message,
		Err:       nil,
	}
}

func Wrap(err error, httpCode int, errorCode, message string) *AppError {
	return &AppError{
		HttpCode:  httpCode,
		ErrorCode: errorCode,
		Message:   message,
		Err:       err,
	}
}

func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		HttpCode:  e.HttpCode,
		ErrorCode: e.ErrorCode,
		Message:   msg,
		Err:       e.Err,
	}
}

func (e *AppError) WithError(err error) *AppError {
	return &AppError{
		HttpCode:  e.HttpCode,
		ErrorCode: e.ErrorCode,
		Message:   e.Message,
		Err:       err,
	}
}
