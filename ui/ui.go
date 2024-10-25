package ui

import (
	"bufio"
	"fmt"
	"os"
	. "sushigo/constants"
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
