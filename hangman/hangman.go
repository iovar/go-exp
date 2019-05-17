package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Println("Welcome to Hangman!")
	startGame()
}

type State int

const (
	Initializing State = 0
	PlayOrQuit         = 1
	Playing      State = 2
	Lost               = 3
	Won                = 4
	Credits            = 5
)

type Game struct {
	word      string
	printWord string
	player    string
	state     State
	found     []string
	failed    []string
	errors    int
	roundLock sync.Mutex
}

func startGame() {
	game := Game{}

	for done := game.Init(); done != false; done = game.PlayOrQuit() {
		go game.Round()
		game.roundLock.Lock()
		game.roundLock.Unlock()
	}
	game.Exit()
}

func (g *Game) Init() bool {
	g.state = Initializing

	fmt.Println("Please enter your name and press ENTER")
	g.player = readLine()

	if len(g.player) == 0 {
		return false
	}

	g.word = getRandomWord()
	g.GetUpdatedWord()
	g.state = Playing

	return true
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	defer func() {
	}()
	return scanner.Text()
}

func readChar() string {
	reader := bufio.NewReader(os.Stdin)
	r, _, _ := reader.ReadRune()
	return string(r)
}

func (g *Game) Round() {
	g.roundLock.Lock()
	fmt.Printf("%s\n", g.printWord)

	c := readChar()
	alreadyUsedFailed := hasLetter(c, g.failed)
	alreadyUsedFound := hasLetter(c, g.found)
	found := g.wordHasLetter(c)

	if !alreadyUsedFailed && !found {
		g.failed = append(g.failed, c)
	} else if !alreadyUsedFound {
		g.found = append(g.found, c)
	}

	if found {
		g.GetUpdatedWord()
	}

	fmt.Printf("Used: %v\n", g.failed)

	g.CheckGameOver(c)
	g.roundLock.Unlock()
}

func (g *Game) wordHasLetter(c string) bool {

	for _, wc := range g.word {
		letter := string(wc)
		if c == letter {
			return true
		}
	}

	return false
}

func (g *Game) CheckGameOver(c string) {
	var remaining int

	remaining = 9 - len(g.failed)

	if remaining <= 0 {
		g.state = Lost
	} else {
		fmt.Printf("You have %d remaining attempts\n", remaining)
	}
}

func (g *Game) GetUpdatedWord() {
	newWord := make([]string, len(g.word))

	for index, c := range g.word {
		letter := string(c)

		if hasLetter(letter, g.found) {
			newWord[index] = letter
		} else {
			newWord[index] = "_"
		}
	}
	g.printWord = strings.Join(newWord, "")
}

func hasLetter(letter string, letters []string) bool {
	for _, c := range letters {
		if strings.EqualFold(letter, c) {
			return true
		}
	}
	return false
}

func (g *Game) Exit() {
	fmt.Println("Thanks for playing! See you again soon!")
}

func (g *Game) PlayOrQuit() bool {
	if g.state == Playing {
		return true
	}

	return false
}

func getRandomWord() string {
	jsonFile, err := os.Open("./mwords.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var words []string

	json.Unmarshal(byteValue, &words)

	return words[rand.Intn(len(words))]
}
