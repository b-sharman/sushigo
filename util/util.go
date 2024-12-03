package util

import (
	"fmt"
	. "sushigo/constants"
)

type Hand [len(QUANTITIES)]int

// A Hand is the collection of cards that are passed between users throughout a
// round. A Board is the collection of cards belonging to a player that is
// scored at the end of each round. Currently, they're stored in the same data
// structure, but it can obviate confusion to distinguish the two.
type Board Hand

func (h Hand) isEmpty() bool {
	for _, count := range h {
		if count > 0 {
			return false
		}
	}
	return true
}

// TODO: make this accept a generic
func PrintHand(hand Hand) {
	for i := 0; i < len(QUANTITIES); i++ {
		if hand[i] != 0 {
			fmt.Printf("%v - %v: %v\n", i, NAMES[i], hand[i])
		}
	}
	fmt.Println()
}

func IsNigiri(ct int) bool {
	return ct == NIGIRI_1 || ct == NIGIRI_2 || ct == NIGIRI_3
}

func IsNigiriOnWasabi(ct int) bool {
	return ct == NIGIRI_1_ON_WASABI || ct == NIGIRI_2_ON_WASABI || ct == NIGIRI_3_ON_WASABI
}

// transform a NIGIRI_n into a NIGIRI_n_ON_WASABI
func Wasabiify(ct int) (int, error) {
	switch ct {
	case NIGIRI_1:
		return NIGIRI_1_ON_WASABI, nil
	case NIGIRI_2:
		return NIGIRI_2_ON_WASABI, nil
	case NIGIRI_3:
		return NIGIRI_3_ON_WASABI, nil
	default:
		return -1, fmt.Errorf("wasabiify received non-nigiri card type %v (expected one of %v, %v, %v)", ct, NIGIRI_1, NIGIRI_2, NIGIRI_3)
	}
}
