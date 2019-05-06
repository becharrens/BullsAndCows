package main

import (
	"BullsAndCows/bullsAndCows"
	"BullsAndCows/mapper"
	"BullsAndCows/possibility"
	"fmt"
	"math/rand"
	"os"
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

	numGuesses := 0
	var bulls, cows int
	var guess string
	var newPoss possibility.Possibility
	var err error
	var possibilitiesFromGuess map[possibility.Possibility]bool
	for {
		numGuesses++
		guess = generateGuessFromPossibility(poss, dgtMapper)

		bulls, cows = bullsAndCows.BullsAndCows(guess, goal)
		if bulls == digits {
			break
		}
		fmt.Printf(
			"Guess: %s - bulls: %d, cows: %d\n", guess, bulls, cows,
		)

		possibilitiesFromGuess = bullsAndCows.GetCandidatesFromResult(
			guess, bulls, cows, values, dgtMapper,
		)

		newPossibilities := make(map[possibility.Possibility]bool)
		for currPoss := range possibilities {
			for guessPoss := range possibilitiesFromGuess {
				newPoss, err = currPoss.Intersect(guessPoss)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				if newPoss != nil {
					newPossibilities[newPoss] = true
				}
			}
		}

		possibilities = newPossibilities
		poss = peekFromMap(possibilities)
	}
	fmt.Printf("Took %d tries to guess the goal: %s", numGuesses, goal)
}

func generateGuessFromPossibility(
	poss possibility.Possibility,
	charMapper mapper.CandidateMapper,
) string {
	guess := ""
	var candidates uint64
	for i := uint(0); i < poss.NumDigits(); i++ {
		candidates = poss.GetDigitCandidates(i)
		guess += generateDgtFromDgtCandidates(candidates, charMapper)
	}
	return guess
}

func generateDgtFromDgtCandidates(
	candidates uint64, charMapper mapper.CandidateMapper,
) string {
	for idx := uint(0); candidates > 0; candidates >>= 1 {
		if candidates&1 > 0 {
			return string(charMapper.MapIdxToChar(idx))
		}
		idx++
	}
	return ""
}

func randomGoal(dgts, values int, charMapper mapper.CandidateMapper) string {
	goal := ""
	var dgt uint
	for i := 0; i < dgts; i++ {
		dgt = uint(rand.Intn(values))
		goal += string(charMapper.MapIdxToChar(dgt))
	}
	return goal
}

func peekFromMap(m map[possibility.Possibility]bool) possibility.Possibility {
	for poss := range m {
		return poss
	}
	return nil
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
