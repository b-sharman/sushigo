package main

import (
	"fmt"
	"log"
	. "sushigo/constants"
	"sushigo/util"
)

type Hand [len(QUANTITIES)]int

func extreme_count(slc []int, comp func(a, b int) int) []int {
	ex_val := slc[0]
	ex_idc := []int{0} // indices of the most extreme values
	for i, num := range slc[1:] {
		i += 1 // correct offset introduced by slicing from 1
		if c := comp(num, ex_val); c > -1 {
			if c == 1 {
				ex_idc = nil
				ex_val = num
			}
			ex_idc = append(ex_idc, i)
		}
	}
	return ex_idc
}

func main() {
	num_players := 2
	cards_per_player := CARD_COUNT - num_players

	hands := make([]Hand, num_players)

	deck := NewDeck()
	for _, hand := range hands {
		cards, err := deck.NextNCards(cards_per_player)
		if err != nil {
			log.Panic(err)
		}
		for _, ct := range cards {
			if util.IsNigiri(ct) && hand[WASABI] > 0 {
				n_on_wasabi, err := util.Wasabiify(ct)
				if err != nil {
					log.Panicf("tried to put non-nigiri card type %v on wasabi", ct)
				}
				hand[n_on_wasabi]++
				hand[WASABI]--
			} else {
				hand[ct]++
			}
		}
		fmt.Printf("Generated hand: %+v\n", hand)
	}

	// this actually scores the hands as they were dealt from the deck
	// in the future the hands have to be distributed to the players (in
	// something like "collections", I guess), and those collections will
	// be stored

	// scores are completely broken until I rewrite the logic

	// scores := score(hands)
	//
	// for _, score := range scores {
	// 	fmt.Println(score)
	// }
}
