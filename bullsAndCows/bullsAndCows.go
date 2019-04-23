package bullsAndCows

import (
	"BullsAndCows/mapper"
	"BullsAndCows/possibility"
)

type Result int

const (
	Bull Result = iota + 1
	Cow
	Empty
)

var (
	possibleResults = []Result{Bull, Cow, Empty}
)

func BullsAndCows(guess, goal string) (int, int) {
	// Assumes both strings are equal in length
	freqTable := buildFreqTable(goal)
	bulls := 0
	cows := 0
	for i := 0; i < len(goal); i++ {
		if guess[i] == goal[i] {
			bulls++
		} else if freqTable[rune(guess[i])] > 0 {
			cows++
		}
	}
	return bulls, cows
}

// Returns the set (map) of all feasible possibilities given the result of the
// previous guess
func GetCandidatesFromResult(
	guess string,
	bulls, cows int,
	numValuesPerDigit uint,
	charMapper mapper.CandidateMapper,
) map[possibility.Possibility]bool {
	freqTable := buildFreqTable(guess)

	possibilities := make(map[possibility.Possibility]bool)
	resultPerm := make(map[rune]Result)

	empty := len(guess) - (bulls + cows)
	result := []int{bulls, cows, empty}

	findCandidatesFromResult(guess, freqTable, result, numValuesPerDigit,
		resultPerm, charMapper, possibilities)
	return possibilities
}

func findCandidatesFromResult(
	guess string,
	freqTable map[rune]int,
	result []int,
	numValuesPerDgt uint,
	resultPerm map[rune]Result,
	charMapper mapper.CandidateMapper,
	possibilities map[possibility.Possibility]bool,
) {
	if len(freqTable) == 0 {
		poss := candidatesForResultPermutation(
			guess, resultPerm, numValuesPerDgt, charMapper,
		)
		possibilities[poss] = true
		return
	}

	char, freq := popFromMap(freqTable)
	for i, res := range possibleResults {
		if result[i] >= freq {
			result[i] -= freq
			resultPerm[char] = res
			findCandidatesFromResult(
				guess, freqTable, result, numValuesPerDgt, resultPerm,
				charMapper, possibilities,
			)
			result[i] += freq
		}
	}
	delete(resultPerm, char)
	freqTable[char] = freq
}

// Given an assignment for the result of a guess, calculate the possibility
// corresponding to it
func candidatesForResultPermutation(
	guess string,
	results map[rune]Result,
	numValuesPerDigit uint,
	charMapper mapper.CandidateMapper,
) possibility.Possibility {
	// Results passed in should be consistent with the guess
	poss := possibility.GetPossibility(uint(len(guess)), numValuesPerDigit)
	var mask, dgtCandidates uint64
	var bitIdx uint
	for i, char := range guess {
		bitIdx = charMapper.MapCharToIdx(char)
		switch results[char] {
		case Empty:
			// Remove char from all entries
			mask = 1 << bitIdx
			for j := uint(0); j < poss.NumDigits(); j++ {
				dgtCandidates = poss.GetDigitCandidates(j)
				poss.SetDigitCandidates(j, dgtCandidates&^mask)
			}
			poss.GetAbsent().SetBit(bitIdx)
		case Cow:
			dgtCandidates = poss.GetDigitCandidates(uint(i))
			dgtCandidates &^= 1 << bitIdx
			poss.SetDigitCandidates(uint(i), dgtCandidates)
			poss.GetPresent().SetBit(bitIdx)
		case Bull:
			poss.SetDigitCandidates(uint(i), 1<<uint(i))
			poss.GetPresent().SetBit(bitIdx)
		}
	}
	return poss
}

func popFromMap(m map[rune]int) (rune, int) {
	for k, v := range m {
		delete(m, k)
		return k, v
	}
	return '0', 0
}

func buildFreqTable(str string) map[rune]int {
	freqTable := make(map[rune]int)
	for _, c := range str {
		freqTable[c]++
	}
	return freqTable
}
