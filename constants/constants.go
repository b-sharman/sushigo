package constants

type Card int

const (
	CHOPSTICKS Card = iota
	DUMPLING
	MAKI_1
	MAKI_2
	MAKI_3
	NIGIRI_1
	NIGIRI_2
	NIGIRI_3
	NIGIRI_1_ON_WASABI
	NIGIRI_2_ON_WASABI
	NIGIRI_3_ON_WASABI
	PUDDING
	SASHIMI
	TEMPURA
	WASABI
)

var (
	QUANTITIES = [...]int{
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

	NAMES = [...]string{
		"Chopsticks",                // CHOPSTICKS
		"Dumpling",                  // DUMPLING
		"Maki 1",                    // MAKI_1
		"Maki 2",                    // MAKI_2
		"Maki 3",                    // MAKI_3
		"Egg Nigiri",                // NIGIRI_1
		"Salmon Nigiri",             // NIGIRI_2
		"Squid Nigiri",              // NIGIRI_3
		"Egg Nigiri (on wasabi)",    // NIGIRI_1_ON_WASABI
		"Salmon Nigiri (on wasabi)", // NIGIRI_2_ON_WASABI
		"Squid Nigiri (on wasabi)",  // NIGIRI_3_ON_WASABI
		"Pudding",                   // PUDDING
		"Sashimi",                   // SASHIMI
		"Tempura",                   // TEMPURA
		"Wasabi",                    // WASABI
	}

	HELPS = [...]string{
		"swap for 2",         // CHOPSTICKS
		"1 3 6 10 15",        // DUMPLING
		"most: 6/3",          // MAKI_1
		"most: 6/3",          // MAKI_2
		"most: 6/3",          // MAKI_3
		"1",                  // NIGIRI_1
		"2",                  // NIGIRI_2
		"3",                  // NIGIRI_3
		"1×3=3",              // NIGIRI_1_ON_WASABI
		"2×3=6",              // NIGIRI_2_ON_WASABI
		"3×3=9",              // NIGIRI_3_ON_WASABI
		"most: 6; least: -6", // PUDDING
		"×3=10",              // SASHIMI
		"×2=5",               // TEMPURA
		"next nigiri ×3",     // WASABI
	}
)

// num of cards per player = CARD_COUNT - numPlayers
const CARD_COUNT = 12

const MIN_PLAYERS = 2
const MAX_PLAYERS = 5

const NUM_ROUNDS = 3

// 1 is passing to the left; -1 is passing to the right
var PASS_DIRECTIONS = [NUM_ROUNDS]int{1, -1, 1}
