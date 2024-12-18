package util

import (
	"fmt"
	. "sushigo/constants"
)

type (
	Hand [len(QUANTITIES)]int

	Board struct {
		data Hand
	}
)

func (board Board) boundsCheck(ct int) error {
	if ct < 0 || ct >= len(board.data) {
		return fmt.Errorf("invalid ct %v", ct)
	}
	return nil
}

// TODO: AddCard should automatically wasabiify
// add 1 card of the corresponding card type to the board
func (board *Board) AddCard(ct int) error {
	err := board.boundsCheck(ct)
	if err != nil {
		return err
	}
	board.data[ct]++
	return nil
}

// add 1 card of the corresponding card types to the board
func (board *Board) AddCards(cts []int) error {
	for _, ct := range cts {
		err := board.AddCard(ct)
		if err != nil {
			return err
		}
	}
	return nil
}

// remove all cards except puddings from the player's board
func (board *Board) Clear() {
	for i := range board.data {
		if i == PUDDING {
			continue
		}
		board.data[i] = 0
	}
}

// return the number of cards corresponding to the given ct on the board
func (board Board) GetQuantity(ct int) (int, error) {
	err := board.boundsCheck(ct)
	if err != nil {
		return -1, err
	}
	return board.data[ct], nil
}

// return the number of cards corresponding to the given ct on the board
// if ct is invalid, panic
func (board Board) GetQuantityNoErr(ct int) int {
	err := board.boundsCheck(ct)
	if err != nil {
		panic("invalid ct passed to GetQuantityNoErr")
	}
	return board.data[ct]
}

// remove 1 card of the corresponding card type from the board
func (board *Board) RemoveCard(ct int) error {
	err := board.boundsCheck(ct)
	if err != nil {
		return err
	}
	if board.data[ct] < 1 {
		return fmt.Errorf("there are no cards of type %v to remove", ct)
	}
	board.data[ct]--
	return nil
}

// set the number of cards corresponding to the given ct on the board
func (board *Board) SetQuantity(ct int, qty int) error {
	err := board.boundsCheck(ct)
	if err != nil {
		return err
	}
	board.data[ct] = qty
	return nil
}

// return a Hand representation of the cards on the board
func (board Board) ToHand() Hand {
	return Hand(board.data)
}

func (h Hand) isEmpty() bool {
	for _, count := range h {
		if count > 0 {
			return false
		}
	}
	return true
}

func PrintHand(hand Hand) {
	for i := range len(QUANTITIES) {
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

// transform a NIGIRI_n_ON_WASABI into a NIGIRI_n
func UnWasabiify(ct int) (int, error) {
	switch ct {
	case NIGIRI_1_ON_WASABI:
		return NIGIRI_1, nil
	case NIGIRI_2_ON_WASABI:
		return NIGIRI_2, nil
	case NIGIRI_3_ON_WASABI:
		return NIGIRI_3, nil
	default:
		return -1, fmt.Errorf("wasabiify received non-nigiri-on-wasabi card type %v (expected one of %v, %v, %v)", ct, NIGIRI_1_ON_WASABI, NIGIRI_2_ON_WASABI, NIGIRI_3_ON_WASABI)
	}
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
