package main

import (
	"fmt"
	"log"
	. "sushigo/constants"
	"sushigo/util"
)

type Hand [len(QUANTITIES)]int

func printHand(hand Hand) {
	for i := 0; i<len(QUANTITIES); i++ {
		fmt.Printf("%v: %v\n", NAMES[i], hand[i])
	}
	fmt.Println()
}

func main() {
	num_players := 2
	cards_per_player := CARD_COUNT - num_players

	hands := make([]Hand, num_players)

	deck := NewDeck()
	for i := range hands {
		cards, err := deck.NextNCards(cards_per_player)
		if err != nil {
			log.Panic(err)
		}
		for _, ct := range cards {
			if util.IsNigiri(ct) && hands[i][WASABI] > 0 {
				n_on_wasabi, err := util.Wasabiify(ct)
				if err != nil {
					log.Panicf("tried to put non-nigiri card type %v (%v) on wasabi", NAMES[ct], ct)
				}
				hands[i][n_on_wasabi]++
				hands[i][WASABI]--
			} else {
				hands[i][ct]++
			}
		}
		printHand(hands[i])
	}

	// this actually scores the hands as they were dealt from the deck
	// in the future the hands have to be distributed to the players (in
	// something like "collections", I guess), and those collections will
	// be stored

	scores := score(hands)

	for _, score := range scores {
		fmt.Println(score)
	}
}
