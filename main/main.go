package main

import (
	"BullsAndCows/bullsAndCows"
	"BullsAndCows/mapper"
	"BullsAndCows/possibility"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	digits = 4
	values = 10
)

func main() {
	rand.Seed(time.Now().UnixNano())
	possibilities := make(map[string]possibility.Possibility)
	poss := possibility.GetPossibility(digits, values)
	possibilities[poss.String()] = poss

	dgtMapper := mapper.GetMapper()
	goal := randomGoal(digits, values, dgtMapper)

	numGuesses := 0
	var bulls, cows int
	var guess string
	var newPoss possibility.Possibility
	var err error
	var possibilitiesFromGuess map[string]possibility.Possibility
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

		newPossibilities := make(map[string]possibility.Possibility)
		fmt.Println(
			"Number of possibility branches to consider:",
			len(possibilities),
		)
		fmt.Println(
			"Number of possibility branches generated by guess:",
			len(possibilitiesFromGuess),
		)
		fmt.Println("")
		for _, currPoss := range possibilities {
			for _, guessPoss := range possibilitiesFromGuess {
				newPoss, err = currPoss.Intersect(guessPoss)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				if newPoss != nil {
					newPossibilities[newPoss.String()] = newPoss
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
		guess += generateDgtFromDgtCandidates(
			candidates, charMapper, poss.BitsPerDigit())
	}
	return guess
}

func generateDgtFromDgtCandidates(
	candidates uint64, charMapper mapper.CandidateMapper, numCandidates uint,
) string {
	nthCandidate := rand.Intn(int(numCandidates))
	numCandidatesSet := 0
	var candIdx int
	candsCopy := candidates
	for {
		candIdx = 0
		for idx := uint(0); candidates > 0; candidates >>= 1 {
			if candidates&1 > 0 {
				numCandidatesSet++
				if candIdx == nthCandidate {
					return string(charMapper.MapIdxToChar(idx))
				}
				candIdx++
			}
			idx++
		}
		candidates = candsCopy
		nthCandidate %= numCandidatesSet
	}
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

func peekFromMap(m map[string]possibility.Possibility) possibility.Possibility {
	for _, poss := range m {
		return poss
	}
	return nil
}
