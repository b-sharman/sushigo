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

	scores := score(hands)

	for _, score := range scores {
		fmt.Println(score)
	}
}
