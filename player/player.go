package player

import (
	"errors"
	. "sushigo/constants"
	"sushigo/ui"
	"sushigo/util"
)

type Player interface {
	// increase the player's count of the given card type
	AddCard(int)

	// remove a card type from the hand and add it to the board
	ChooseCard(*util.Hand) (int, error)

	// return the cards that the player has played this round
	GetBoard() util.Board

	// decrement the number of wasabis on the player's board
	// useful when converting a nigiri to a nigiri_on_wasabi
	RemoveWasabi()
}

type HumanPlayer struct {
	board util.Board
}

func (hp *HumanPlayer) AddCard(ct int) {
	hp.board[ct]++
}

func (hp *HumanPlayer) ChooseCard(hand *util.Hand) (int, error) {
	return ui.GetCardType(&hp.board, hand), nil
}

func (hp HumanPlayer) GetBoard() util.Board {
	return hp.board
}

func (hp *HumanPlayer) RemoveWasabi() {
	(*hp).board[WASABI]--
}

type ComputerPlayer struct {
	board util.Board
}

func (cp *ComputerPlayer) AddCard(ct int) {
	cp.board[ct]++
}

func (cp *ComputerPlayer) ChooseCard(hand *util.Hand) (int, error) {
	// for now, just pick the first valid card
	for ct, count := range *hand {
		if count > 0 {
			return ct, nil
		}
	}
	return -1, errors.New("hand has no cards")
}

func (cp ComputerPlayer) GetBoard() util.Board {
	return cp.board
}

func (cp *ComputerPlayer) RemoveWasabi() {
	(*cp).board[WASABI]--
}
