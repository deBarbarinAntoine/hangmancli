package HangmanCLI

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Save struct {
	Name       string
	Score      int
	Word       string
	Difficulty string
	Dictionary string
}

type Parameters struct {
	Name           string
	DictionaryPath string
	Difficulty     int
}

var savedGames []Save

// saveGame saves the current game in fileName.
func saveGame(fileName string) {
	currentGame := Save{
		Name:       MyGame.Name,
		Score:      MyGame.Score,
		Word:       MyGame.Word,
		Difficulty: toStringDifficulty(MyGame.Difficulty),
		Dictionary: dictionaryName(MyGame.Dictionary),
	}
	newEntry, err := json.Marshal(currentGame)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	newEntry = append([]byte{',', '\n'}, newEntry...)
	_, err = file.Write(newEntry)
	if err != nil {
		log.Fatal(err)
	}
}

// retreiveSavedGames retreive all saved games present in fileName and put it in savedEntries.
func retreiveSavedGames(fileName string) {
	savedEntries, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucune sauvegarde détectée...", CLEARCOLOR)
		return
	}
	savedEntries = append([]byte{'[', '\n'}, savedEntries...)
	savedEntries = append(savedEntries, '\n', ']')
	err = json.Unmarshal(savedEntries, &savedGames)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orangered), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	}
}

// chargeParameters retreive the parameters present in fileName and changes all corresponding variables.
func chargeParameters(fileName string) {
	var savedParameters Parameters
	savedEntries, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucun fichier de configuration détecté...", CLEARCOLOR)
		time.Sleep(time.Second * 1)
		return
	}
	err = json.Unmarshal(savedEntries, &savedParameters)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "Il est conseillé de supprimer le fichier config.txt\n    dans le dossier Files afin de résoudre le problème.", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orange), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	} else {
		MyGame.Name = savedParameters.Name
		MyGame.Dictionary = savedParameters.DictionaryPath
		MyGame.Difficulty = savedParameters.Difficulty
	}
}

// saveParameters saves all current parameters in fileName for later use.
func saveParameters(fileName string) {
	currentParameters := Parameters{
		Name:           MyGame.Name,
		DictionaryPath: MyGame.Dictionary,
		Difficulty:     MyGame.Difficulty,
	}
	newEntry, err := json.Marshal(currentParameters)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fileName, newEntry, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
