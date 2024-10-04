package main

const CHOPSTICKS = 0
const DUMPLING = 1
const MAKI_1 = 2
const MAKI_2 = 3
const MAKI_3 = 4
const NIGIRI_1 = 5
const NIGIRI_2 = 6
const NIGIRI_3 = 7
const NIGIRI_1_ON_WASABI = 8
const NIGIRI_2_ON_WASABI = 9
const NIGIRI_3_ON_WASABI = 10
const PUDDING = 11
const SASHIMI = 12
const TEMPURA = 13
const WASABI = 14

var QUANTITIES = [...]int{
	4,  // CHOPSTICKS
	14, // DUMPLING
	6,  // MAKI_1
	12, // MAKI_2
	8,  // MAKI_3
	5,  // NIGIRI_1
	10, // NIGIRI_2
	5,  // NIGIRI_3
	0,  // NIGIRI_1_ON_WASABI
	0,  // NIGIRI_2_ON_WASABI
	0,  // NIGIRI_3_ON_WASABI
	10, // PUDDING
	14, // SASHIMI
	14, // TEMPURA
	6,  // WASABI
}

// num of cards per player = CARD_COUNT - num_players
const CARD_COUNT = 12
