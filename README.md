# My-Words-App
I build this tool to better manage my English vocabulary study. 

Tha app is built with `Go`, `sqlite`, and these awesome cli frameworks / packages:
- [cobra](https://github.com/spf13/cobra)
- [promptui](https://github.com/manifoldco/promptui)
- [go-pretty](https://github.com/jedib0t/go-pretty)

## Main features
- CRUD operations for words and definitions. Support cli input and csv file upload.
- Search a word's definition through [dictionary api](https://dictionaryapi.dev/).
- Take quizzes and review the quiz results.
- Review the incorrect words from the last quiz and do a quiz again on only those words.
- `-h` or `--help` on commands to see exact usages.

## Commands
Usage: `my-words-app [command]`

Available Commands:

  `delete`:      Delete a word from database
  
  `help`:        Help about any command
  
  `init`:        Initialize the app. **Run this if it's your first time using the app!!!**
  
  `quiz`:        Take a quiz of [count] words from database. You need minimum 10 words in database to run this command.
  
  `search`:      Search word definition from the dictionary api. After search, you can add the word to the database.
  
  `select`:      Print or create a csv file of the words and definitions in database. If no flags are provided, will select latest 10 words.
  
  transcript  Show transcript of [quizId] quiz, if no quizId is provided, will show the latest quiz. Use --history to get the quizId.
  
  `upsert`:      Insert or update a word and its definition to the database. Use --file flag to upsert with a csv file.

## Usage
1. Install [Go](https://go.dev/) (version >= 1.18) if you haven't already
2. Clone this repository
3. Add a `.env` file to the root. Refer to `.env.example` 

    The app uses [this](https://apilayer.com/marketplace/dymt-api) api to provide word suggestion of there's a spelling misktake when using `search` command. You can either get your own api key, or disable it by setting `ENABLE_WORD_SUGGESTION=false` in `.env`
    
4. Run `go build`
5. `--help` to see the help 
``` 
./my-words-app --help
```