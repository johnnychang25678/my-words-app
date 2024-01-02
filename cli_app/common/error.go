package common

import (
	"fmt"
)

type ErrorCode string

const (
	UpsertError       ErrorCode = "UPSERT_ERROR"
	OpenFileError     ErrorCode = "OPEN_FILE_ERROR"
	CreateFileError   ErrorCode = "CREATE_FILE_ERROR"
	ReadCsvError      ErrorCode = "READ_CSV_ERROR"
	DbError           ErrorCode = "DB_ERROR"
	SelectError       ErrorCode = "SELECT_ERROR"
	QuizError         ErrorCode = "QUIZ_ERROR"
	TranscriptError   ErrorCode = "TRANSCRIPT_ERROR"
	DeleteError       ErrorCode = "DELETE_ERROR"
	SearchError       ErrorCode = "SEARCH_ERROR"
	SearchNoDefError  ErrorCode = "SEARCH_NO_DEF_ERROR"
	ApiError          ErrorCode = "API_ERROR"
	InvalidInputError ErrorCode = "INVALID_INPUT_ERROR"
	UnknownError      ErrorCode = "UNKNOWN_ERROR"
)

type AppError struct {
	ErrorCode ErrorCode
	Err       error
}

func (a AppError) PrintAppError() {
	fmt.Println("**ERROR**", fmt.Sprintf("code: %s, message: %s", a.ErrorCode, a.Err.Error()))
}
