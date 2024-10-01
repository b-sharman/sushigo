package main

import "fmt"

type Nigiri struct {
	value     uint
	on_wasabi bool
}

type Hand struct {
	dumpling, pudding, sashimi, tempura uint
	maki     []uint
	nigiri   []Nigiri
}

func main() {
	hands := []Hand{
		{
			0,                               // dumpling
			2,                               // pudding
			0,                               // sashimi
			3,                               // tempura
			[]uint{3, 3},                    // maki
			[]Nigiri{{2, true}, {3, false}}, // nigiri
		},
		{
			3,                    // dumpling
			0,                    // pudding
			2,                    // sashimi
			3,                    // tempura
			[]uint{2, 3, 2},      // maki
			[]Nigiri{{2, false}}, // nigiri
		},
	}

	scores := []uint{0, 0}

	// types of cards that don't depend on other players
	for i := range scores {
		scores[i] += []uint{0, 1, 3, 6, 10, 15}[hands[i].dumpling]
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

	// // types of cards that depend on other players
	// for i := range scores {
	//        // award 6 points to the player with the most puddings
	//        // award 6 points to the player with the most puddings
    //        // handle ties for both
	//        f
	// }

	for _, score := range scores {
        fmt.Println(score)
    }
}
