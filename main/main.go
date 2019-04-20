package main

import (
	"BullsAndCows/possibility"
	"fmt"
	"math/rand"
)

func main() {

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
