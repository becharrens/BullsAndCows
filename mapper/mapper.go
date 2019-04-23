package mapper

type CandidateMapper interface {
	MapCharToIdx(char rune) uint
	MapIdxToChar(char uint) rune
}

type DigitMapper struct {
}

// Maps the index of a bit within a number (bitmap) to its corresponding
// character
func (mapper *DigitMapper) MapIdxToChar(idx uint) rune {
	return rune('0' + idx)
}

// Maps a character to its index within a bitmap
func (mapper *DigitMapper) MapCharToIdx(c rune) uint {
	// Assumes character is a digit
	return uint(c - '0')
}

func GetMapper() CandidateMapper {
	return &DigitMapper{}
}
