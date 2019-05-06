package possibility

import (
	"bytes"
	"errors"
	"fmt"
)

const maxBitsPerDigit = 64

// For now assume that number of values per digit cannot exceed 64

type Possibility interface {
	Intersect(other Possibility) (Possibility, error)
	NumEntries() int
	NumDigits() uint
	BitsPerDigit() uint
	GetEntry(idx uint) uint64
	SetEntry(idx uint, value uint64) bool
	GetDigitCandidates(dgt uint) uint64
	SetDigitCandidates(dgt uint, value uint64) bool
	GetAbsent() BitMap
	GetPresent() BitMap
	String() string
}

type PossibleNums struct {
	possibleNums   []uint64
	valuesPerDigit uint
	digits         uint
	presentDigits  BitMap
	absentDigits   BitMap
}

func (poss *PossibleNums) NumEntries() int {
	return len(poss.possibleNums)
}

// Sets an entry in the possibleNums array to the given value
func (poss *PossibleNums) SetEntry(idx uint, entry uint64) bool {
	if idx >= uint(len(poss.possibleNums)) {
		return false
	}
	poss.possibleNums[idx] = entry
	return true
}

// Retrieves the ith element in the array
func (poss *PossibleNums) GetEntry(idx uint) uint64 {
	return poss.possibleNums[idx]
}

// Returns the number of values a digit in the number can take
func (poss *PossibleNums) BitsPerDigit() uint {
	return poss.valuesPerDigit
}

// Returns the number of digits in the number
func (poss *PossibleNums) NumDigits() uint {
	return poss.digits
}

// Returns a uint64 representing the bitmap for the candidates
// at that particular index in the number.
// Note: The digit at index 0 represents the most representative
// digit in the number
func (poss *PossibleNums) GetDigitCandidates(dgt uint) uint64 {
	if dgt >= poss.NumDigits() {
		return 0
	}
	startBit := poss.BitsPerDigit() * dgt
	entry := startBit / maxBitsPerDigit
	offset := startBit % maxBitsPerDigit
	mask := uint64(1<<poss.BitsPerDigit() - 1)
	result := (poss.GetEntry(entry) >> offset) & mask
	numBits := maxBitsPerDigit - offset
	if numBits < poss.BitsPerDigit() {
		mask = uint64(1<<(poss.BitsPerDigit()-numBits) - 1)
		result += (poss.GetEntry(entry+1) & mask) << numBits
	}
	return result
}

// Sets the digit candidates for the digit at position dgt within the number
// (index 0 at leftmost digit in the number). As the number of values a digit
// can take need not be a multiple of 64, the digit candidates may be split
// across two entries
func (poss *PossibleNums) SetDigitCandidates(dgt uint, value uint64) bool {
	if dgt >= poss.NumDigits() {
		return false
	}
	startBit := poss.BitsPerDigit() * dgt
	entryIdx := startBit / 64
	offset := startBit % 64
	entry := poss.GetEntry(entryIdx)
	remBits := maxBitsPerDigit - offset
	if remBits >= poss.BitsPerDigit() {
		entry = setBitsInEntry(value, entry, poss.BitsPerDigit(), offset)
		poss.SetEntry(entryIdx, entry)
	} else {
		entry = setBitsInEntry(value, entry, remBits, offset)
		value >>= remBits
		poss.SetEntry(entryIdx, entry)
		entry = poss.GetEntry(entryIdx + 1)
		entry = setBitsInEntry(
			value, entry, poss.BitsPerDigit()-remBits, 0)
		poss.SetEntry(entryIdx+1, entry)
	}
	return true
}

// Calculates the intersection of two possibilities. It calculates the
// intersection of the candidates for each digit in the number, and if any
// of the resulting candidate values are inconsistent (intersection = 0),
// return nil, an empty intersection
func (poss *PossibleNums) Intersect(other Possibility) (Possibility, error) {
	if other.NumDigits() != poss.NumDigits() ||
		other.BitsPerDigit() != poss.BitsPerDigit() {
		return nil, errors.New("possibilities don't match")
	}

	if !poss.GetAbsent().And(other.GetPresent()).IsEmpty() ||
		!poss.GetPresent().And(other.GetAbsent()).IsEmpty() {
		return nil, nil
	}

	possibility := &PossibleNums{
		valuesPerDigit: poss.valuesPerDigit,
		digits:         poss.digits,
		possibleNums:   make([]uint64, len(poss.possibleNums)),
		absentDigits:   poss.GetAbsent().Or(other.GetAbsent()),
		presentDigits:  poss.GetPresent().Or(other.GetPresent()),
	}

	for i := 0; i < len(poss.possibleNums); i++ {
		possibility.possibleNums[i] =
			poss.GetEntry(uint(i)) & other.GetEntry(uint(i))
	}

	entryIdx := uint(0)
	offset := uint(0)
	empty := false
	for digitsChecked := uint(0); digitsChecked < poss.NumDigits(); digitsChecked++ {
		empty, entryIdx, offset =
			possibility.areDigitCandidatesEmpty(
				entryIdx, offset, poss.BitsPerDigit())
		if empty {
			return nil, nil
		}
	}
	return possibility, nil
}

// Checks if the candidates for a particular digit are empty. It handles
// the case where candidates are split between two entries using the offset
// and numBits parameters, returning the entry index and the offset of the start
// of the digit's candidates
func (poss *PossibleNums) areDigitCandidatesEmpty(
	entryIdx, offset, numBits uint,
) (bool, uint, uint) {
	entry := poss.GetEntry(entryIdx)
	entry >>= offset
	if allZeros(entry, numBits) {
		if maxBitsPerDigit-offset >= poss.BitsPerDigit() {
			return true, 0, 0
		}
		return poss.areDigitCandidatesEmpty(entryIdx+1, 0,
			poss.BitsPerDigit()-(maxBitsPerDigit-offset))
	}
	if maxBitsPerDigit-offset <= poss.BitsPerDigit() {
		return false, entryIdx + 1,
			poss.BitsPerDigit() - (maxBitsPerDigit - offset)
	}
	return false, entryIdx, offset + poss.BitsPerDigit()
}

// String representation of a possibility, printing the candidates for each
// digit as a binary string
func (poss PossibleNums) String() string {
	var result bytes.Buffer
	formatString := "%" + fmt.Sprintf("%d", poss.BitsPerDigit()) + "b "
	for i := uint(0); i < poss.NumDigits(); i++ {
		result.Write(
			[]byte(fmt.Sprintf(formatString, poss.GetDigitCandidates(i))))
	}
	return result.String()
}

// Sets numBits number of bits in a uint64 bit number, starting at
// an offset 'offset'
func setBitsInEntry(value, entry uint64, numBits, offset uint) uint64 {
	mask := uint64(1<<numBits - 1)
	value = (value & mask) << offset
	mask <<= offset
	entry = entry &^ mask
	return entry | value
}

// Checks if the first numbits digits in a number are all 0 (starting from the
// right)
func allZeros(num uint64, numBits uint) bool {
	return num&uint64(1<<numBits-1) == 0
}

func (poss *PossibleNums) GetAbsent() BitMap {
	return poss.absentDigits
}

func (poss *PossibleNums) GetPresent() BitMap {
	return poss.presentDigits
}

// Constructs a possibility with all the fields are set properly and all the
// digit candidates are initialised to 1s
func GetPossibility(numDigits uint, valuesPerDigit uint) Possibility {
	if numDigits <= 0 || valuesPerDigit == 0 {
		return nil
	}

	possibility := &PossibleNums{
		valuesPerDigit: valuesPerDigit, digits: numDigits,
		absentDigits:  GetBitMap(0, numDigits),
		presentDigits: GetBitMap(0, numDigits),
	}
	numBits := numDigits * valuesPerDigit
	numEntries := numBits / 64
	if numBits%64 > 0 {
		numEntries++
	}
	possibleDgts := make([]uint64, numEntries)
	for i := uint(0); i < numEntries-1; i++ {
		possibleDgts[i] = 1<<maxBitsPerDigit - 1
	}
	possibleDgts[numEntries-1] = 1<<(numBits%64) - 1
	possibility.possibleNums = possibleDgts

	return possibility
}
