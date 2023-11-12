package HangmanCLI

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type Game struct {
	Name                     string
	Score                    int
	Word                     string
	Difficulty               int
	Dictionary               string
	WordDisplay, RunesPlayed []rune
	words                    []string
	nbLettersFound, nbErrors int
}

var (
	hangman []string
	myGame  Game
)

const (
	ALREADYPLAYED  = 1
	CORRECTRUNE    = 2
	INCORRECTRUNE  = 3
	CORRECTWORD    = 5
	INCORRECTWORD  = 6
	INCORRECTINPUT = 8

	EASY      = 4
	MEDIUM    = 7
	DIFFICULT = 10
	LEGENDARY = 13

	ONGOING = 14
	WIN     = 15
	LOOSE   = 16
)

// initGame initialises the game with user data and all variables necessary.
func initGame(name, dictionary string, difficulty int) {
	retreiveWords(dictionary)
	myGame.Word = string(chooseWord(difficulty))
	myGame.WordDisplay = []rune(strings.Repeat("_ ", len(myGame.Word)))
	myGame.WordDisplay = hint()
}

// checkInputFormat checks if the input's format is right and returns the result with a boolean.
func checkInputFormat(input string) bool {
	for _, char := range input {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return false
		}
	}
	return true
}

// checkWord checks if the player found the word. Returns true if he found it and false if not.
func checkWord(try string) int {
	if try == myGame.Word {
		return CORRECTWORD
	}
	return INCORRECTWORD
}

// Function that changes the wordDisplay to replace the '_' character with the rune played if it is in the word.
func displayWord(char rune) []rune {
	for i, r := range myGame.Word {
		if r == char {
			myGame.WordDisplay[i*2] = char - 32
			myGame.nbLettersFound++
			myGame.Score += 10
		}
	}
	return myGame.WordDisplay
}

// revealWord reveals all runes in wordDisplay.
func revealWord() []rune {
	for i, r := range myGame.Word {
		myGame.WordDisplay[i*2] = r - 32
	}
	return myGame.WordDisplay
}

// nbRemainingLetters returns the number of letters still not found in the word.
func nbRemainingLetters() int {
	var result int
	for _, char := range myGame.WordDisplay {
		if char == '_' {
			result++
		}
	}
	return result
}

// retreiveWords retreive the words from the selected dictionary.
func retreiveWords(dictionary string) {
	if dictionary == "" {
		dictionary = "../Files/Dictionaries/ods5.txt"
	}
	content, err := os.ReadFile(dictionary)
	if err != nil {
		log.Fatal(err)
	}
	if checkDictionary() {
		myGame.words = strings.Split(string(content), "\n")
	} else {
		fmt.Println(colorCode(Red), "Erreur d'acquisition des mots du dictionnaire", CLEARCOLOR)
		time.Sleep(time.Second * 2)
	}
}

// retreiveHangman retreives the hangman in /Files/hangman.txt and stores it in hangman.
func retreiveHangman() {
	hangman = append(hangman[0:0])
	content, err := os.ReadFile("../Files/hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	var line int
	var str string
	for _, char := range content {
		str += string(char)
		if char == '\n' {
			line++
		}
		if line == 8 {
			hangman = append(hangman, str)
			str = ""
			line = 0
		}
	}
}

// checkRune checks if the rune played is already played, correct or incorrect.
func checkRune(char rune) int {
	for _, r := range myGame.RunesPlayed {
		if r == char {
			return ALREADYPLAYED
		}
	}
	for _, r := range strings.ToUpper(myGame.Word) {
		if r == char {
			myGame.RunesPlayed = append(myGame.RunesPlayed, char)
			return CORRECTRUNE
		}
	}
	myGame.RunesPlayed = append(myGame.RunesPlayed, char)
	return INCORRECTRUNE
}

func countScore(result int) {
	switch result {
	case ALREADYPLAYED:
		break
	case CORRECTRUNE:
		myGame.Score += 10
	case INCORRECTRUNE:
		myGame.Score -= 5
		myGame.nbErrors++
	case CORRECTWORD:
		myGame.Score += 11 * nbRemainingLetters()
	case INCORRECTWORD:
		myGame.nbErrors += 2
		myGame.Score -= 5
	}
}

// clearGameData clears the previous' game's data to start a new one.
func clearGameData() {
	myGame.Score = 0
	myGame.Word = ""
	myGame.WordDisplay = append(myGame.WordDisplay[0:0])
	myGame.RunesPlayed = append(myGame.RunesPlayed[0:0])
	myGame.words = append(myGame.words[0:0])
	myGame.nbLettersFound = 0
	myGame.nbErrors = 0
}

// hint reveal a random rune in wordDisplay.
func hint() []rune {
	if myGame.Difficulty != LEGENDARY {
		i := rand.Intn(len(myGame.Word) - 1)
		char := []rune(myGame.Word)[i]
		myGame.WordDisplay[i*2] = char - 32
	}
	return myGame.WordDisplay
}

// chooseWord chooses randomly a word from words (the dictionary's words' list) according to the difficulty set previously.
func chooseWord(difficulty int) string {
	var possibleWords []string
	for _, str := range myGame.words {
		if len(str) >= difficulty-2 && len(str) <= difficulty {
			possibleWords = append(possibleWords, str)
		}
		if difficulty == LEGENDARY {
			if len(str) > difficulty {
				possibleWords = append(possibleWords, str)
			}
		}
	}
	if len(possibleWords) < 10 {
		var i int
		for _, str := range myGame.words {
			i++
			if len(str) == difficulty-i-2 || len(str) == difficulty+i {
				possibleWords = append(possibleWords, str)
				if len(possibleWords) > 10 {
					break
				}
			}
		}
	}
	return possibleWords[rand.Intn(len(possibleWords)-1)]
}

// dictionaryName returns the name of the dictionary.
func dictionaryName(dictionary string) string {
	switch dictionary {
	case "../Files/Dictionaries/ods5.txt":
		return "Scrabble"
	case "../Files/Dictionaries/ospd3_expurgated.txt":
		return "Anglais"
	case "../Files/Dictionaries/italiano.txt":
		return "Italien"
	default:
		return "Personnalisé"
	}
}

// sortTopTenGames sort the saved games by score in decreasing order.
func sortTopTenGames() []Save {
	sort.SliceStable(savedGames, func(i, j int) bool { return savedGames[i].Score > savedGames[j].Score })
	return savedGames
}

// toStringDifficulty returns the name of the difficulty.
func toStringDifficulty(difficulty int) string {
	switch difficulty {
	case EASY:
		return "Facile"
	case MEDIUM:
		return "Intermédiaire"
	case DIFFICULT:
		return "Difficile"
	case LEGENDARY:
		return "Légendaire"
	default:
		return "Inconnu"
	}
}

// checkDictionary checks if the dictionary is usable or not and changes the case.
func checkDictionary() bool {
	for i, str := range myGame.words {
		myGame.words[i] = strings.ToLower(str)
		str = myGame.words[i]
		for _, char := range str {
			if char < 'a' || char > 'z' {
				return false
			}
		}
	}
	return true
}

// checkEndGame verify if the game is finished and returns the status and true if the game is finished or false is the game is still ongoing.
func checkEndGame() (int, bool) {
	if myGame.nbErrors >= len(hangman)-1 {
		return LOOSE, true
	}
	if strings.Join(strings.Split(string(myGame.WordDisplay), " "), "") == strings.ToUpper(myGame.Word) {
		return WIN, true
	}
	return ONGOING, false
}

// Run executes the game, obviously... :)
func Run() {
	chargeParameters("../Files/config.txt")
	principalMenu()
}