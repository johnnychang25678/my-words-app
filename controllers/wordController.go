package controllers

import "github.com/johnnychang25678/my_words_app_2.0/errors"

type Word struct {
	Id         int    `json:"id"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type IWordController interface {
	GetWords() ([]Word, *errors.AppError)
}

// should satisfy IWordController
type WordController struct {
}

func (w WordController) GetWords() ([]Word, *errors.AppError) {
	return []Word{{1, "johnny", "name"}}, nil
}

func NewWordController() IWordController {
	return &WordController{}
}
