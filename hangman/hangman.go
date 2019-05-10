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
	player    string
	state     State
	letters   []string
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
	g.PrintWord()
	c := readChar()

	if !hasLetter(c, g.letters) {
		g.letters = append(g.letters, c)
	}

	fmt.Printf("You entered: %v\n", g.letters)
	g.roundLock.Unlock()
}

func (g *Game) PrintWord() {
	for _, c := range g.word {
		letter := string(c)
		if hasLetter(letter, g.letters) {
			fmt.Print(letter)
		} else {
			fmt.Print("_")
		}
	}
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
