package plr

import (
	. "sushigo/constants"
	"sushigo/util"
)

type Player struct {
	Board   util.Board
	Chooser Reasoner
}

type Reasoner interface {
	// round number, player index, boards, hand
	ChooseCard(int, int, []util.Board, util.Hand) ([]Card, error)
}

func BoardsFromPlayers(players []*Player) []util.Board {
	boards := make([]util.Board, 0, len(players))
	for _, player := range players {
		boards = append(boards, player.Board)
	}
	return boards
}
