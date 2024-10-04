package main

import (
	"errors"
	"math/rand/v2"
	"strconv"
)

type Deck []int

func new_deck() Deck {
	deck := make(Deck, len(QUANTITIES))
	copy(deck, QUANTITIES[:])
	return deck
}

// return a slice of n indices x_i where NAMES[x_i] is the type of card
func (deck *Deck) next_n_cards(n int) ([]int, error) {
	remaining_cards := 0
	for _, q := range *deck {
		remaining_cards += q
	}
	if remaining_cards < n {
		return nil, errors.New("deck has " + strconv.Itoa(remaining_cards) + " cards but " + strconv.Itoa(n) + " cards were requested")
	}

	card_types := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var ct int
		// keep generating a random card type until the deck has a card of that type
		// not the most efficient but the simplest I can think of now
		for ct = rand.IntN(len(QUANTITIES)); (*deck)[ct] == 0; ct = rand.IntN(len(QUANTITIES)) {
		}
		(*deck)[ct]--
		card_types = append(card_types, ct)
	}

	return card_types, nil
}
