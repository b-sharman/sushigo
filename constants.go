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
	"pudding",
	"sashimi",
	"tempura",
	"wasabi",
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
	10, // pudding
	14, // sashimi
	14, // tempura
	6,  // wasabi
}

// num of cards per player = CARD_COUNT - num_players
const CARD_COUNT = 12
