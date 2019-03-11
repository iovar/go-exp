package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	Player      = "P"
	Ghost       = "@"
	Empty       = "_"
	LevelLength = 10
)

func main() {
	level := make([]string, LevelLength)
	resetLevel(level)

	var mutex sync.Mutex
	fmt.Print("Let's go!\n\n")
	runGameLoop(level, &mutex)
}

func resetLevel(level []string) {
	pPos, gPos := getStartPositions()
	for index, _ := range level {
		if index == pPos {
			level[index] = Player
		} else if index == gPos {
			level[index] = Ghost
		} else {
			level[index] = Empty
		}
	}
}

func getStartPositions() (int, int) {
	pPos := rand.Intn(LevelLength)
	gPos := rand.Intn(LevelLength)
	if gPos == pPos {
		return getStartPositions()
	}
	return gPos, pPos
}

func getMovement() int {
	return rand.Intn(2)
}

func getPos(level []string, entity string) int {
	for index, value := range level {
		if value == entity {
			return index
		}
	}
	return -1
}

func getNextPos(current int, offset int) int {
	return (current + offset) % LevelLength
}

func performRound(level []string) []string {
	defer func() {
		if recover() != nil {
			fmt.Print(level)
		}
	}()
	userPos := getPos(level, Player)
	userMove := getMovement()
	nextUserPos := getNextPos(userPos, userMove)

	ghostPos := getPos(level, Ghost)
	ghostMove := getMovement()
	nextGhostPos := getNextPos(ghostPos, ghostMove)

	level[userPos] = Empty
	level[ghostPos] = Empty
	level[nextUserPos] = Player
	level[nextGhostPos] = Ghost

	if nextUserPos == nextGhostPos {
		checkEndGame(userPos, ghostPos)
		time.Sleep(3 * time.Second)
		resetLevel(level)
	}

	fmt.Print("\r")
	fmt.Print(level)
	return level
}

func checkEndGame(userPos int, ghostPos int) {
	if userPos < ghostPos {
		fmt.Print("\nWell done!\n\n")
	} else if userPos > ghostPos {
		fmt.Print("\nYou are dead!\n\n")
	}
}

func runGameLoop(level []string, mutex *sync.Mutex) {
	mutex.Lock()
	for {
		performRound(level)
		time.Sleep(500 * time.Millisecond)
	}
}
