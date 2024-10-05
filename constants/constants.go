package constants

const (
	CHOPSTICKS = iota
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

var NAMES = [...]string{
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

// num of cards per player = CARD_COUNT - num_players
const CARD_COUNT = 12
