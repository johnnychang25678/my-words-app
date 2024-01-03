package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/johnnychang25678/my_words_app_2.0/controllers"
	"github.com/johnnychang25678/my_words_app_2.0/errors"
	"github.com/valyala/fasthttp"
)

// **************** response signatures *********************
type response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func newOkResponse(data any) response {
	return response{Message: "ok", Data: data}
}

func newErrorResponse(errorMessage string) response {
	return response{Message: errorMessage, Data: nil}
}

// helper function
func handlerHelper(c *fasthttp.RequestCtx, statusCode int, jsonBytes []byte) {
	c.SetStatusCode(statusCode)
	c.SetContentType("application/json")
	c.Write(jsonBytes)
}

func errorHandler(c *fasthttp.RequestCtx, e *errors.AppError) {
	j, err := json.Marshal(newErrorResponse(e.Error()))
	if err != nil {
		fmt.Println("json marshal error:")
		fmt.Println(err)
	} else {
		handlerHelper(c, e.StatusCode, j)
	}
}

func okHandler(c *fasthttp.RequestCtx, data any) {
	j, err := json.Marshal(newOkResponse(data))
	if err != nil {
		fmt.Println("json marshal error:")
		fmt.Println(err)
	} else {
		handlerHelper(c, 200, j)
	}
}

// **************** manage dependency *********************
type App struct {
	wordController controllers.IWordController
}

func initApp() App {
	return App{
		wordController: controllers.NewWordController(),
	}
}

// **********************  handlers ***************************
func (app App) getWordsHandler(c *fasthttp.RequestCtx) {
	words, err := app.wordController.GetWords()
	if err != nil {
		errorHandler(c, err)
		return
	}
	okHandler(c, words)
}

func testHandler(c *fasthttp.RequestCtx) {
	okHandler(c, []int{1, 2, 3})
}

func main() {
	app := initApp()

	// routing
	server := &fasthttp.Server{
		Handler: func(c *fasthttp.RequestCtx) {
			path := string(c.Path())
			switch path {
			case "/":
				testHandler(c)
			case "/words":
				app.getWordsHandler(c)
			}
		},
	}

	port := ":8080"
	fmt.Println("server running at ", port)
	if err := server.ListenAndServe(port); err != nil {
		log.Fatal(err)
	}
}
