# My-Words-App
I build this command-line app to better manage my English vocabulary study.

Live Demo: 
<iframe  title="YouTube video player" width="480" height="390" src="https://www.youtube.com/watch?v=wJMVpKq_YYM&ab_channel=Chia-MingChang" frameborder="0" allowfullscreen></iframe>

### What is the problem you're solving?

In the past, this was the process when I was trying to memorize vocabularies:

See a word I don't understand -> google it -> Add the word and its definition to an excel file -> Review and try to quiz myself by hiding the definition column on excel -> Repeat

As you can see, this process is tedious and hard to keep track of my progress. For example, if I want to take a quiz only on the incorrect words from last quiz, it's quite a manual process using Excel. 

The app is a centralized place for me to do word searching, studying and quizzing.

### Why command-line app? Why not mobile or web?

Because I am a nerd and I love cli! I feel like I can better concentrate with cli and just using my keyboard. Plus I've always wanted to build a cli tool with Go - this is the first one I've ever built.
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
  
  `transcript`:  Show transcript of [quizId] quiz, if no quizId is provided, will show the latest quiz. Use --history to get the quizId.
  
  `upsert`:      Insert or update a word and its definition to the database. Use --file flag to upsert with a csv file.

## Usage
You can simply go to [/example](https://github.com/johnnychang25678/my-words-app/tree/main/example) and download the app to start using it! More details in README in example.

Or you can build the project following these steps:
1. Install [Go](https://go.dev/) (version >= 1.18) if you haven't already
2. Clone this repository
3. Add a `.env` file to the root. Refer to `.env.example` 

    The app uses [this](https://apilayer.com/marketplace/dymt-api) api to provide word suggestion if there's a spelling misktake when using `search` command. You can either get your own api key, or disable it by setting `ENABLE_WORD_SUGGESTION=false` in `.env`
    
4. Run `go build`
5. `--help` to see the help 
``` 
./my-words-app --help
```
## Tech stack
The app is built with `Go`, `sqlite`, and these awesome cli frameworks / packages:
- [cobra](https://github.com/spf13/cobra)
- [promptui](https://github.com/manifoldco/promptui)
- [go-pretty](https://github.com/jedib0t/go-pretty)

## Future Improvements
- Add more colors and styles.
- Support different languages.
- Search for synonyms when adding words. Option to add synonyms to database.
- Search and add the word's example sentences.
- ... 