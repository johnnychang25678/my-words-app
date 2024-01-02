# Example
1. Download `my-words-app` and `example_words.csv` to the same directory
2. See help message for all the commands
```
./my-words-app --help
```
3. If you do not have a `.env` file in the same directory, the app will create one.
4. Initalize app
```
./my-words-app init
```
5. Add 100 words to the database
```
./my-words-app upsert -f example_words.csv
```
6. Select the words
```
./my-words-app select --all 
```
7. Take a quiz
```
./my-words-app quiz
```
8. Take a quiz again only on the incorrect words from last quiz
```
./my-words-app quiz --incorrect
```
9. See quiz records
```
./my-words-app transcript --history
```
10. See certain transcript by provding quiz id
```
./my-words-app transcript 1
```
11. Search for word definition and add to database
```
./my-words-app search nebulous
```
12. Select the word you just added with -c flag
```
./my-words-app select -c 1
```
13. If you have the api key of the word suggestion api, you can try miss typing a word when search
```
./my-words-app search scissorr
```