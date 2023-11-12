package HangmanCLI

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// runCmd executes the command and arguments put in the parameters.
func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// clearTerminal clears the terminal using the corresponding command.
func clearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func principalMenu() {
	for {
		switch menu("------- MENU PRINCIPAL -------", "Nouvelle partie", "Dictionnaire", "Meilleurs scores", "Quitter") {
		case "Nouvelle partie":
			setGame()
		case "Dictionnaire":
			changeDictionary()
		case "Meilleurs scores":
			topScores()
		case "Quitter":
			saveParameters("../Files/config.txt")
			os.Exit(0)
		}
	}
}

func setGame() {
	var incorrectInput bool
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "------- INITIALISATION DU JEU -------", CLEARCOLOR)
		fmt.Println()
		if incorrectInput {
			fmt.Println(colorCode(Red), "Nom saisi incorrect (entre 3 et 15 caractères, sans nombres ni signes)", CLEARCOLOR)
		}
		fmt.Print(colorCode(Forestgreen), "Saisissez votre nom : ", colorCode(Aquamarine))
		if nameInput() {
			break
		} else {
			incorrectInput = true
		}
	}
	setDifficulty()
}

func setDifficulty() {
	switch menu("------- DIFFICULTÉ -------", "Facile", "Intermédiaire", "Difficile", "Légendaire") {
	case "Facile":
		MyGame.Difficulty = EASY
	case "Intermédiaire":
		MyGame.Difficulty = MEDIUM
	case "Difficile":
		MyGame.Difficulty = DIFFICULT
	case "Légendaire":
		MyGame.Difficulty = LEGENDARY
	}
	play()
}

func play() {
	retreiveHangman()
	initGame(MyGame.Name, MyGame.Dictionary, MyGame.Difficulty)
	var status int
	var gameHasEnded bool
	var previousResult int
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "------------------------- HANGMAN -------------------------", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Forestgreen), "Nom : ", colorCode(Aquamarine), MyGame.Name, colorCode(Forestgreen), "\tDifficulté : ", colorCode(Aquamarine), toStringDifficulty(MyGame.Difficulty), colorCode(Forestgreen), "\tDictionnaire : ", colorCode(Aquamarine), dictionaryName(MyGame.Dictionary), colorCode(Forestgreen), "\tScore : ", colorCode(Aquamarine), MyGame.Score, CLEARCOLOR)
		fmt.Println()
		fmt.Println(hangman[MyGame.nbErrors])
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), string(MyGame.WordDisplay), CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Forestgreen), "Lettres déjà jouées : ", colorCode(Orange), string(MyGame.RunesPlayed), CLEARCOLOR)
		fmt.Println()
		if status, gameHasEnded = checkEndGame(); gameHasEnded {
			time.Sleep(time.Second * 2)
			break
		}
		switch previousResult {
		case ALREADYPLAYED:
			fmt.Println(colorCode(Orangered), "Cette lettre a déjà été jouée !", CLEARCOLOR)
			fmt.Println()
		case INCORRECTINPUT:
			fmt.Println(colorCode(Orangered), "Saisie invalide !", CLEARCOLOR)
			fmt.Println()
		}
		fmt.Print(colorCode(Deepskyblue), "Proposez une lettre ou un mot : ", colorCode(Aquamarine))
		previousResult = input()
	}
	endGame(status)
}

func endGame(status int) {
	for {
		clearTerminal()
		if status == WIN {
			fmt.Println(colorCode(Cyan), "\tFÉLICITATIONS, VOUS AVEZ GAGNÉ !", CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Le mot était ", strings.ToUpper(MyGame.Word), CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Votre score est : ", MyGame.Score, CLEARCOLOR)
		} else if status == LOOSE {
			fmt.Println(colorCode(Orange), "\tGAME OVER !", CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Red), hangman[MyGame.nbErrors], CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Le mot était ", MyGame.Word, CLEARCOLOR)
		}

		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "0.  Retour au menu principal")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "0" {
			break
		}
	}
	saveGame("../Files/scores.txt")
	clearGameData()
}

func changeDictionary() {
	switch menu("------- CHANGER DE DICTIONNAIRE -------", "Scrabble français", "Scabble Anglais", "Italien", "Retour") {
	case "Scrabble français":
		MyGame.Dictionary = "../Files/Dictionaries/ods5.txt"
	case "Scabble Anglais":
		MyGame.Dictionary = "../Files/Dictionaries/ospd3_expurgated.txt"
	case "Italien":
		MyGame.Dictionary = "../Files/Dictionaries/italiano.txt"
	case "Retour":
		break
	}
}

func topScores() {
	retreiveSavedGames("../Files/scores.txt")
	sortTopTenGames()
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "-------------------- MEILLEURS SCORES --------------------", CLEARCOLOR)
		for i, game := range savedGames {
			fmt.Println(colorCode(Forestgreen), i+1, "\t", game.Name, "\t", game.Score, "\t", game.Difficulty, "\t", game.Dictionary, "\t", CLEARCOLOR)
			if i > 10 {
				break
			}
		}
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "0.  Retour au menu principal")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "0" {
			break
		}
	}
}
