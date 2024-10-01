package main

import "fmt"

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
			1,                    // pudding
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
	puddings := make([]int, 0, len(hands))
	for _, hand := range hands {
		puddings = append(puddings, hand.pudding)
	}
	// TODO: use a sort function instead
	least_p_val := puddings[0]
	most_p_val := puddings[0]
	least_p_idx := []int{0}
	most_p_idx := []int{0}
	for i, num := range puddings[1:] {
		// correct offset introduced by slicing puddings
		i += 1
		if num <= least_p_val {
			if num < least_p_val {
				least_p_idx = nil
				least_p_val = num
			}
			least_p_idx = append(least_p_idx, i)
		}
		if num >= most_p_val {
			if num > most_p_val {
				most_p_idx = nil
				most_p_val = num
			}
			most_p_idx = append(most_p_idx, i)
		}
	}
	if len(hands) > 2 {
		for _, idx := range least_p_idx {
			scores[idx] -= 6 / len(least_p_idx)
		}
	}
	for _, idx := range most_p_idx {
		scores[idx] += 6 / len(most_p_idx)
	}

	for _, score := range scores {
		fmt.Println(score)
	}
}
