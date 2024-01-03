package errors

type AppError struct {
	StatusCode   int
	ErrorMessage string
}

func (ae AppError) Error() string {
	return ae.ErrorMessage
}
