package common

import (
	"fmt"
)

type ErrorCode string

const (
	UpsertError     ErrorCode = "UPSERT_ERROR"
	OpenFileError   ErrorCode = "OPEN_FILE_ERROR"
	ReadCsvError    ErrorCode = "READ_CSV_ERROR"
	DbError         ErrorCode = "DB_ERROR"
	SelectError     ErrorCode = "SELECT_ERROR"
	QuizError       ErrorCode = "QUIZ_ERROR"
	TrasncriptError ErrorCode = "TRANSCRIPT_ERROR"
)

type AppError struct {
	ErrorCode ErrorCode
	Err       error
}

func (a AppError) PrintAppError() {
	fmt.Println("**ERROR**", fmt.Sprintf("code: %s, message: %s", a.ErrorCode, a.Err.Error()))
}
