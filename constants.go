package main

var NAMES = [...]string{
	"chopsticks",
	"dumpling",
	"maki_1",
	"maki_2",
	"maki_3",
	"nigiri_1",
	"nigiri_2",
	"nigiri_3",
	"nigiri_1_on_wasabi",
	"nigiri_2_on_wasabi",
	"nigiri_3_on_wasabi",
	"pudding",
	"sashimi",
	"tempura",
	"wasabi",
}

var NAMES_LOOKUP = map[string]int{
	"chopsticks": 0,
	"dumpling": 1,
	"maki_1": 2,
	"maki_2": 3,
	"maki_3": 4,
	"nigiri_1": 5,
	"nigiri_2": 6,
	"nigiri_3": 7,
	"nigiri_1_on_wasabi": 8,
	"nigiri_2_on_wasabi": 9,
	"nigiri_3_on_wasabi": 10,
	"pudding": 11,
	"sashimi": 12,
	"tempura": 13,
	"wasabi": 14,
}

var QUANTITIES = [...]int{
	4,  // chopsticks
	14, // dumpling
	6,  // maki_1
	12, // maki_2
	8,  // maki_3
	5,  // nigiri_1
	10, // nigiri_2
	5,  // nigiri_3
	0,  // nigiri_1_on_wasabi
	0, // nigiri_2_on_wasabi
	0,  // nigiri_3_on_wasabi
	10, // pudding
	14, // sashimi
	14, // tempura
	6,  // wasabi
}

// num of cards per player = CARD_COUNT - num_players
const CARD_COUNT = 12
