package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Nigiri struct {
	value     int
	on_wasabi bool
}

type Hand struct {
	dumpling, pudding, sashimi, tempura int
	maki                                []int
	nigiri                              []Nigiri
}

func totalMakis(hand Hand) int {
	var sum int = 0
	for _, val := range hand.maki {
		sum += val
	}
	return sum
}

func main() {
	hands := []Hand{
		{
			0,                               // dumpling
			2,                               // pudding
			0,                               // sashimi
			3,                               // tempura
			[]int{3, 3},                     // maki
			[]Nigiri{{2, true}, {3, false}}, // nigiri
		},
		{
			3,                    // dumpling
			2,                    // pudding
			2,                    // sashimi
			3,                    // tempura
			[]int{2, 3, 2},       // maki
			[]Nigiri{{2, false}}, // nigiri
		},
	}

	scores := make([]int, len(hands), len(hands))

	// types of cards that don't depend on other players
	for i := range scores {
		scores[i] += []int{0, 1, 3, 6, 10, 15}[hands[i].dumpling]
		for _, n := range hands[i].nigiri {
			if n.on_wasabi {
				scores[i] += 3 * n.value
			} else {
				scores[i] += n.value
			}
		}
		scores[i] += 5 * (hands[i].sashimi / 3)
		scores[i] += 5 * (hands[i].tempura / 2)
	}

	// types of cards that depend on other players

	// award 6 points to the player with the most puddings
	// award -6 points to the player with the least puddings
	// handle ties for both
	puddings := make([][2]int, 0, len(hands))
	for i, hand := range hands {
		puddings = append(puddings, [2]int{i, hand.pudding})
	}
	slices.SortFunc(puddings, func(a, b [2]int) int {
		return cmp.Compare(b[1], a[1])
	})
	if len(hands) > 2 {
		// penalize players with the least puddings
		min_puddings := puddings[len(puddings)-1][1]
		// indices of hands that have the least puddings
		var to_lower []int
		for i := len(puddings) - 1; i >= 0 && puddings[i][1] == min_puddings; i-- {
			to_lower = append(to_lower, puddings[i][0])
		}
		for _, idx := range to_lower {
			scores[idx] -= 6 / len(to_lower)
		}
	}
	// reward players with the most puddings
	max_puddings := puddings[0][1]
	// indices of hands that have the most puddings
	var to_raise []int
	for i := 0; i < len(puddings) && puddings[i][1] == max_puddings; i++ {
		to_raise = append(to_raise, puddings[i][0])
	}
	for _, idx := range to_raise {
		scores[idx] += 6 / len(to_raise)
	}

	for _, score := range scores {
		fmt.Println(score)
	}
}
