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
		num_results, err := fmt.Scanln(&num)
		// num_results, err := fmt.Scanf("%d\n", &num)
		if err == nil && num_results == 1 && num >= MIN_PLAYERS && num <= MAX_PLAYERS {
			valid = true
		}
		if num_results != 1 {
			stdin.ReadString('\n') // clear stdin
		}
	}
	return num
}

func GetCardType(board *util.Board, hand *util.Hand) int {
	fmt.Println("\nThe hand you're holding has:")
	util.PrintHand(*hand)

	stdin := bufio.NewReader(os.Stdin)
	var ct int
	valid := false
	for !valid {
		fmt.Print("Enter the number corresponding to the card you'd like to play: ")
		num_results, err := fmt.Scanln(&ct)
		if err == nil && ct >= 0 && ct < len(QUANTITIES) && hand[ct] > 0 {
			valid = true
		}
		if num_results != 1 {
			stdin.ReadString('\n') // clear stdin
		}
	}
	return ct
}
