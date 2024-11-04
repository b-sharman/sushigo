package plr

import (
	"errors"
	. "sushigo/constants"
	"sushigo/ui"
	"sushigo/util"
)

type Player struct {
	board   util.Board
	IsHuman bool
}

// increase the player's count of the given card type
func AddCard(player *Player, ct int) {
	player.board[ct]++
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

func ChooseCard(cp *Player, hand *util.Hand) ([]int, error) {
	if cp.IsHuman {
		return humanChooseCard(cp, hand)
	}
	return computerChooseCard(cp, hand)
}

// return the cards that the player has played this round
func GetBoard(player *Player) util.Board {
	return player.board
}

// decrement the number of chopsticks on the player's board
func RemoveChopsticks(player *Player) error {
	if player.board[CHOPSTICKS] < 1 {
		return errors.New("there are no chopsticks to remove")
	}
	player.board[CHOPSTICKS]--
	return nil
}

// decrement the number of wasabis on the player's board
func RemoveWasabi(player *Player) {
	player.board[WASABI]--
}

// remove a card type from the hand and add it to the board
func humanChooseCard(hp *Player, hand *util.Hand) ([]int, error) {
	hasChopsticks := hp.board[CHOPSTICKS] > 0
	return ui.GetCardType(hasChopsticks, hand), nil
}

func computerChooseCard(cp *Player, hand *util.Hand) ([]int, error) {
	// for now, just pick the first valid card
	for ct, count := range *hand {
		if count > 0 {
			return []int{ct}, nil
		}
	}
	return []int{}, errors.New("hand has no cards")
}
