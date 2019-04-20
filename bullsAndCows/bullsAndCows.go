package bullsAndCows

import (
	"BullsAndCows/possibility"
)

type Result int

const (
	Bull Result = iota + 1
	Cow
	Empty
)

func getCandidatesFromResult(
	guess string,
	bulls, cows,
	numValuesPerDigit int,
) map[possibility.Possibility]bool {
	return nil
}

func candidatesForResultPermutation(
	guess string,
	results map[rune]Result,
	numValuesPerDigit uint,
) possibility.Possibility {
	// Results passed in should be consistent with the guess
	poss := possibility.GetPossibility(uint(len(guess)), numValuesPerDigit)
	for i, char := range guess {
		switch results[char] {
		case Empty:
			// Remove char from all entries
		case Bull:
			// poss.SetDigitCandidates(i, 1<<i)
		}
	}
	return poss
}
