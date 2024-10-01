package main

import "fmt"

// ooo I'm sure this is so idiomatic
const CHOPSTICKS = 7
const DUMPLING   = 2
const MAKI       = 3
const NIGIRI     = 4
const PUDDING    = 5
const SASHIMI    = 1
const TEMPURA    = 0
const WASABI     = 6

func main() {
    hand := map[int][]int {
        CHOPSTICKS: {
		},
        DUMPLING: {
		},
        MAKI: {
		},
        NIGIRI: {
            2, 3
		},
        PUDDING: {
		},
        SASHIMI: {
		},
        TEMPURA: {
            0,0,0
		},
        WASABI: {
            0
		},
    }

	var wn map[*int]*int // maps wasabi to nigiri
	wn[&hand[WASABI][0]] = &hand[NIGIRI][0]
	for _, nigiri := range wn {
	    nigiri.value *= 3
	}

    score := 0
    // for _, value := range DUMPLING {
    // }
    // for _, value := range MAKI {
    // }
    for _, value := range NIGIRI {
        score += value
    }
    // for _, value := range PUDDING {
    // }
    // for _, value := range SASHIMI {
    // }
    for _, value := range TEMPURA {
    }
    for _, value := range WASABI {
    }

}
