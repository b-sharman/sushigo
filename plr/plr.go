package plr

import (
	"fmt"
	. "sushigo/constants"
	"sushigo/util"
)

type Player struct {
	board   util.Board
	Chooser Reasoner
}

type Reasoner interface {
	ChooseCard(int, int, []util.Board, util.Hand) ([]int, error)
}

// increase the player's count of the given card type
func AddCard(player *Player, ct int) {
	player.board[ct]++
}

func BoardsFromPlayers(players []*Player) []util.Board {
	boards := make([]util.Board, 0, len(players))
	for _, player := range players {
		boards = append(boards, GetBoard(player))
	}
	return boards
}

// remove all cards except puddings from the player's board
func ClearBoard(player *Player) {
	for i := range player.board {
		if i == PUDDING {
			continue
		}
		player.board[i] = 0
	}
}

// return the cards that the player has played this round
func GetBoard(player *Player) util.Board {
	return player.board
}

// decrement the number of chopsticks on the player's board
func RemoveCard(player *Player, ct int) error {
	if ct < 0 || ct > len(QUANTITIES) {
		return fmt.Errorf("invalid card type %v", ct)
	}
	if player.board[ct] < 1 {
		return fmt.Errorf("there are no cards of type %v to remove", ct)
	}
	player.board[ct]--
	return nil
}
