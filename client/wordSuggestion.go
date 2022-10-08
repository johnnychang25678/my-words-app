package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/johnnychang25678/my-words-app/common"
)

type WordSuggestionReponse struct {
	IsModified   bool   `json:"is_modified"`
	OriginalText string `json:"original_text"`
	Result       string `json:"result"`
}

type WordSuggestionError struct {
	Message string `json:"message"`
}

type wordSuggestionClient struct {
	httpClient *resty.Client
}

func NewWordSuggestionClient() wordSuggestionClient {
	return wordSuggestionClient{
		httpClient: resty.New().
			SetBaseURL("https://api.apilayer.com/dymt/did_you_mean_this").
			SetHeader("apikey", os.Getenv("WORD_SUGGESTION_API_KEY")),
	}
}

func (w wordSuggestionClient) GetWordSuggestion(wrongWord string) (string, *common.AppError) {
	var res WordSuggestionReponse
	var apiErr WordSuggestionError
	resp, err := w.httpClient.R().
		SetQueryParam("q", wrongWord).
		SetResult(&res).
		SetError(&apiErr).
		Get("")
	if err != nil {
		return "", &common.AppError{ErrorCode: common.ApiError, Err: err}
	}
	if resp.StatusCode() != 200 {
		if apiErr.Message != "" {
			return "", &common.AppError{ErrorCode: common.ApiError, Err: fmt.Errorf(apiErr.Message)}
		}
		fmt.Println(resp)
		return "", &common.AppError{ErrorCode: common.UnknownError, Err: fmt.Errorf("unknown error")}
	}
	return res.Result, nil
}
