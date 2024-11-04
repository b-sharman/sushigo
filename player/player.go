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
	ChooseCard(*util.Hand) ([]int, error)

	// remove all cards except puddings from the player's board
	ClearBoard()

	// return the cards that the player has played this round
	GetBoard() util.Board

	// decrement the number of chopsticks on the player's board
	RemoveChopsticks() error

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

func (hp *HumanPlayer) ChooseCard(hand *util.Hand) ([]int, error) {
	hasChopsticks := hp.board[CHOPSTICKS] > 0
	return ui.GetCardType(hasChopsticks, hand), nil
}

func (hp *HumanPlayer) ClearBoard() {
	for i := range (*hp).board {
		if i == PUDDING {
			continue
		}
		(*hp).board[i] = 0
	}
}

func (hp HumanPlayer) GetBoard() util.Board {
	return hp.board
}

func (hp *HumanPlayer) RemoveChopsticks() error {
	if (*hp).board[CHOPSTICKS] < 1 {
		return errors.New("there are no chopsticks to remove")
	}
	(*hp).board[CHOPSTICKS]--
	return nil
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

func (cp *ComputerPlayer) ChooseCard(hand *util.Hand) ([]int, error) {
	// for now, just pick the first valid card
	for ct, count := range *hand {
		if count > 0 {
			return []int{ct}, nil
		}
	}
	return []int{}, errors.New("hand has no cards")
}

func (cp *ComputerPlayer) ClearBoard() {
	for i := range (*cp).board {
		if i == PUDDING {
			continue
		}
		(*cp).board[i] = 0
	}
}

func (cp ComputerPlayer) GetBoard() util.Board {
	return cp.board
}

func (cp *ComputerPlayer) RemoveChopsticks() error {
	if (*cp).board[CHOPSTICKS] < 1 {
		return errors.New("there are no chopsticks to remove")
	}
	(*cp).board[CHOPSTICKS]--
	return nil
}

func (cp *ComputerPlayer) RemoveWasabi() {
	(*cp).board[WASABI]--
}
