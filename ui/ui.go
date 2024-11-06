package ui

import (
	"bufio"
	"fmt"
	"os"
	. "sushigo/constants"
	"sushigo/util"
)

func GetNumPlayers() int {
	stdin := bufio.NewReader(os.Stdin)
	var num int
	valid := false
	for !valid {
		fmt.Print("Enter the number of players: ")
		numResults, err := fmt.Scanln(&num)
		if err == nil && numResults == 1 && num >= MIN_PLAYERS && num <= MAX_PLAYERS {
			valid = true
		}
		if numResults != 1 {
			stdin.ReadString('\n') // clear stdin
		}
	}
	return num
}

// returns an int either corresponding to a card type or -1 if the player wishes to use chopsticks
func getSingleCardType(hand util.Hand, canBeChopsticks bool) int {
	stdin := bufio.NewReader(os.Stdin)
	var ct int
	valid := false
	for !valid {
		fmt.Print("Enter the number corresponding to the card you'd like to play")
		if canBeChopsticks {
			fmt.Print(", or -1 to use chopsticks")
		}
		fmt.Print(": ")
		numResults, err := fmt.Scanln(&ct)
		if err == nil && ct < len(QUANTITIES) && ((ct >= 0 && hand[ct] > 0) || (canBeChopsticks && ct == -1)) {
			valid = true
		}
		if numResults != 1 {
			stdin.ReadString('\n') // clear stdin
		}
	}
	return ct
}

func GetCardType(hasChopsticks bool, hand util.Hand) []int {
	fmt.Println("\nThe hand you're holding has:")
	util.PrintHand(hand)

	// chopsticks can only be used when there are at least two cards in the hand
	canBeChopsticks := false
	if hasChopsticks {
		total := 0
		for _, count := range hand {
			total += count
			if total >= 2 {
				canBeChopsticks = true
				break
			}
		}
	}

	ct := getSingleCardType(hand, canBeChopsticks)
	if ct == -1 {
		// TODO: if the player chooses a nigiri first and wasabi
		// afterwards, prompt them to ensure that they intentionally
		// selected the cards in that order.
		return []int{
			getSingleCardType(hand, false),
			getSingleCardType(hand, false),
		}
	}
	return []int{ct}
}

func PrintScores(scores [][]int, numPlayers int, roundIdx int) {
	fmt.Print("\t\t")
	for j := 0; j < numPlayers; j++ {
		fmt.Printf("Player %v\t", j)
	}
	fmt.Print("\n\t\t")
	for j := 0; j < numPlayers; j++ {
		fmt.Print("--------\t")
	}
	fmt.Println()
	for j := 0; j <= roundIdx; j++ {
		fmt.Printf("Round %v\t\t", j)
		for _, score := range scores[j] {
			fmt.Printf("%v\t\t", score)
		}
		fmt.Println()
	}
	if roundIdx == NUM_ROUNDS-1 {
		fmt.Print("\t\t")
		for j := 0; j < numPlayers; j++ {
			fmt.Print("--------\t")
		}
		fmt.Print("\nTotal\t\t")
		for j := 0; j < numPlayers; j++ {
			score := 0
			for _, roundScores := range scores {
				score += roundScores[j]
			}
			fmt.Printf("%v\t\t", score)
		}
		fmt.Println()
	}
}
