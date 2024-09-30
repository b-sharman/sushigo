package main

import "fmt"

// type Card interface {
// 	Score([]Instance) int
// }

type (
	Tempura struct{}

	Dumpling struct{}

	Nigiri struct {
		value     int // TODO: confine to 1, 2, or 3
		on_wasabi bool
	}
)

/*
each type of card has a struct
each of those structs implements Scorable
Scorable is an interface offering a Score method accepting an array/slice of structs and returning an int
To calculate the total score of a hand, iterate through all the possible types of cards; for each type, generate a list of the cards in that hand with that type and pass it to the type's Score().

problem: the score of some types of cards (maki rolls, pudding) depends on other players' cards
*/

func main() {
	{
		hand := [3]Nigiri{
			{1, false},
			{2, false},
			{3, true},
		}
		sum := 0
		for _, m := range hand {
			if m.on_wasabi {
				sum += 3 * m.value
			} else {
				sum += m.value
			}
		}
		fmt.Println("Nigiri:", sum)
	}

	{
		hand := [5]Tempura{
			{},
			{},
			{},
			{},
			{},
		}
		sum := 5 * (len(hand) / 2)
		fmt.Println("Tempura:", sum)
	}

	{
		dumpling_scores := [6]int{0, 1, 3, 6, 10, 15}
		hand := [3]Dumpling{
			{},
			{},
			{},
		}
		sum := dumpling_scores[min(4, len(hand))]
		fmt.Println("Dumpling:", sum)
	}
}
