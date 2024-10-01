package main

import (
	"cmp"
	"fmt"
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
	hands := []Hand{
		{
			0,                               // dumpling
			0,                               // pudding
			0,                               // sashimi
			3,                               // tempura
			[]int{},                         // maki
			[]Nigiri{{2, true}, {3, false}}, // nigiri
		},
		{
			3,                    // dumpling
			0,                    // pudding
			2,                    // sashimi
			3,                    // tempura
			[]int{1},             // maki
			[]Nigiri{{2, false}}, // nigiri
		},
		{
			3,                    // dumpling
			0,                    // pudding
			2,                    // sashimi
			3,                    // tempura
			[]int{},              // maki
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

	// puddings
	puddings := make([]int, 0, len(hands))
	for _, hand := range hands {
		puddings = append(puddings, hand.pudding)
	}
	// penalize players with the least puddings
	if len(hands) > 2 {
		least_p_idc := extreme_count(puddings, func(a, b int) int { return -1 * cmp.Compare(a, b) })
		for _, idx := range least_p_idc {
			scores[idx] -= 6 / len(least_p_idc)
		}
	}
	// award players with the most puddings
	most_p_idc := extreme_count(puddings, cmp.Compare)
	for _, idx := range most_p_idc {
		scores[idx] += 6 / len(most_p_idc)
	}

	// makis
	makis := make([]int, 0, len(hands))
	for _, hand := range hands {
		makis = append(makis, totalMakis(hand))
	}
	// award players with the most makis
	most_m_val := -1
	sndmost_m_val := -1
	var most_m_idc []int    // indices of the most extreme values
	var sndmost_m_idc []int // indices of the second-most extreme values
	for i, num := range makis {
		if num >= most_m_val {
			if num > most_m_val {
				// the former most is the new sndmost
				sndmost_m_val = most_m_val
				sndmost_m_idc = make([]int, len(most_m_idc))
				copy(sndmost_m_idc, most_m_idc)
				most_m_idc = nil
				most_m_val = num
			}
			most_m_idc = append(most_m_idc, i)
		} else if num >= sndmost_m_val {
			if num > sndmost_m_val {
				sndmost_m_val = num
				sndmost_m_idc = nil
			}
			sndmost_m_idc = append(sndmost_m_idc, i)
		}
	}
	if most_m_val > 0 {
		for _, idx := range most_m_idc {
			scores[idx] += 6 / len(most_m_idc)
		}
	}
	if len(most_m_idc) < 2 && sndmost_m_val > 0 {
		for _, idx := range sndmost_m_idc {
			scores[idx] += 3 / len(sndmost_m_idc)
		}
	}

	for _, score := range scores {
		fmt.Println(score)
	}
}
