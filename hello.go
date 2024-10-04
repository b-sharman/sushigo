package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Nigiri struct {
	value     int
	on_wasabi bool
}

type Hand struct {
	dumpling, pudding, sashimi, tempura, wasabi int
	maki                                        []int
	nigiri                                      []Nigiri
}

func totalMakis(hand Hand) int {
	var sum int = 0
	for _, val := range hand.maki {
		sum += val
	}
	return sum
}

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

	deck := new_deck()
	// for i := 0; i<num_players; i++ {
	for i, hand := range hands {
		cards, err := deck.next_n_cards(cards_per_player)
		if err != nil {
			fmt.Println("Uh oh, the deck ran out of cards! This probably happened because there are " + strconv.Itoa(num_players) + " players. The maximum allowed is 5.")
		}
		for _, ct := range cards {
			fmt.Printf("Handling ct %v (%v)\n", ct, NAMES[ct])
			switch name, _, _ := strings.Cut(NAMES[ct], "_"); name {
			case "chopsticks":
			case "dumpling":
				hand.dumpling++
			case "maki_1":
				hand.maki = append(hand.maki, 1)
			case "maki_2":
				hand.maki = append(hand.maki, 2)
			case "maki_3":
				hand.maki = append(hand.maki, 3)
			case "nigiri_1":
				hand.nigiri = append(hand.nigiri, Nigiri{1, hand.wasabi > len(hand.nigiri)})
			case "nigiri_2":
				hand.nigiri = append(hand.nigiri, Nigiri{2, hand.wasabi > len(hand.nigiri)})
			case "nigiri_3":
				hand.nigiri = append(hand.nigiri, Nigiri{3, hand.wasabi > len(hand.nigiri)})
			case "pudding":
				hand.pudding++
			case "sashimi":
				hand.sashimi++
			case "tempura":
				hand.tempura++
			case "wasabi":
				hand.wasabi++
			}
			fmt.Printf("Handled ct %v (%v)\n", ct, NAMES[ct])
		}
		fmt.Printf("Generated hand: %+v\n", hand)
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
