package main

import "fmt"

type Nigiri struct {
	value     int
	on_wasabi bool
}

type Hand struct {
	dumpling int
	maki     []int
	nigiri   []Nigiri
	pudding  int
	sashimi  int
	tempura  int
}

func main() {
	hand := Hand{
		0,                               // dumpling
		[]int{},                         // maki
		[]Nigiri{{2, true}, {3, false}}, // nigiri
		0,                               // pudding
		0,                               // sashimi
		3,                               // tempura
	}

	score := 0
	score += []int{0, 1, 3, 6, 10, 15}[hand.dumpling]
	for _, n := range hand.nigiri {
		if n.on_wasabi {
			score += 3 * n.value
		} else {
			score += n.value
		}
	}
	score += 5 * (hand.sashimi / 3)
	score += 5 * (hand.tempura / 2)

	fmt.Println(score)
}
