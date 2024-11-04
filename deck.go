package main

import (
	"fmt"
	"math/rand/v2"
	. "sushigo/constants"
)

type Deck []int

func NewDeck() Deck {
	deck := make(Deck, len(QUANTITIES))
	copy(deck, QUANTITIES[:])
	return deck
}

// return a slice of n indices x_i where NAMES[x_i] is the type of card
func (deck *Deck) NextNCards(n int) ([]int, error) {
	remainingCards := 0
	for _, q := range *deck {
		remainingCards += q
	}
	if remainingCards < n {
		return nil, fmt.Errorf("deck has %v cards but %v cards were requested", remainingCards, n)
	}

	cardTypes := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var ct int
		// keep generating a random card type until the deck has a card of that type
		// not the most efficient but the simplest I can think of now
		for ct = rand.IntN(len(QUANTITIES)); (*deck)[ct] == 0; ct = rand.IntN(len(QUANTITIES)) {
		}
		(*deck)[ct]--
		cardTypes = append(cardTypes, ct)
	}

	return cardTypes, nil
}
