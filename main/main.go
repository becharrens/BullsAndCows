package main

import (
	"BullsAndCows/bullsAndCows"
	"BullsAndCows/mapper"
	"BullsAndCows/possibility"
	"fmt"
	"math/rand"
)

const (
	digits = 4
	values = 10
)

func main() {
	possibilities := make(map[possibility.Possibility]bool)
	poss := possibility.GetPossibility(digits, values)
	possibilities[poss] = true

	dgtMapper := mapper.GetMapper()

	goal := randomGoal(digits, values, dgtMapper)
	var bulls, cows int
	var guess string
	for {
		guess = generateGuessFromPossibility(poss, dgtMapper)
		fmt.Printf(
			"Guess: %s - bulls: %d, cows: %d\n", guess, bulls, cows,
		)
		bulls, cows = bullsAndCows.BullsAndCows(guess, goal)
		if bulls == digits {
			break
		}
		// TODO: Finish main loop
	}
}

func generateGuessFromPossibility(
	poss possibility.Possibility,
	charMapper mapper.CandidateMapper,
) string {
	return ""
}

func randomGoal(dgts, values int, charMapper mapper.CandidateMapper) string {
	return ""
}

func test() {
	poss := possibility.GetPossibility(10, 10)
	newPoss := possibility.GetPossibility(10, 10)
	for j := 0; j < 3; j++ {
		for i := uint(0); i < poss.NumDigits(); i++ {
			newPoss.SetDigitCandidates(i, uint64(rand.Int63n(1<<63-1)))
		}
		fmt.Println(poss)
		fmt.Println(newPoss)
		intersect, _ := poss.Intersect(newPoss)
		fmt.Println(intersect, "\n")
		poss = intersect
	}
}
