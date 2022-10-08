package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DictionaryResponse []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string        `json:"definition"`
			Example    string        `json:"example"`
			Synonyms   []interface{} `json:"synonyms"`
			Antonyms   []interface{} `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

type DictionaryError struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Resolution string `json:"resolution"`
}

type dictionaryClient struct {
	httpClient *resty.Client
}

func NewDictionaryClient() dictionaryClient {
	return dictionaryClient{
		httpClient: resty.New().SetBaseURL("https://api.dictionaryapi.dev/api/v2/entries/en/"),
	}
}

func (d dictionaryClient) GetDefinitions(word string) ([]string, error) {
	var dicRes DictionaryResponse
	var dicErr DictionaryError
	resp, err := d.httpClient.R().
		SetPathParam("word", word).
		SetResult(&dicRes).
		SetError(&dicErr).
		Get("/{word}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		if dicErr.Title != "" {
			return nil, fmt.Errorf(dicErr.Title)
		}
		return nil, fmt.Errorf("unknown error")
	}

	var output []string
	definitions := dicRes[0].Meanings[0].Definitions
	for _, d := range definitions {
		output = append(output, d.Definition)
	}
	return output, nil
}
