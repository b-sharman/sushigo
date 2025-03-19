// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"math/rand"

	. "sushigo/constants"
	"sushigo/util"
)

type Computer struct{}

/* choose a card for this turn given a hand and some other information
 *
 * arguments:
 * roundNum: 0, 1, or 2, corresponding to which round it is (used to determine passing direction of cards)
 * myIdx: boards[myIdx] = my board
 * boards: slice of all players' boards, not including the cards they have chosen this round
 * hand: the hand of cards I can choose from
 *
 * returns:
 * []int: slice of card type(s) chosen
 * error: incorrect parameters, or bug in ChooseCard implementation
 */
func (cp *Computer) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]Card, error) {
	// randomly pick a card for now
	numCards := 0
	for _, qty := range hand {
		numCards += qty
	}
	seed := rand.Intn(numCards)
	numLookedAt := 0
	for ct, qty := range hand {
		numLookedAt += qty
		if numLookedAt >= seed {
			return []Card{Card(ct)}, nil
		}
	}
	return nil, errors.New("something went wrong when picking a card")
}
