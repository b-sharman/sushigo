// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	. "sushigo/constants"
	"sushigo/score"
	"sushigo/util"
)

/* some nice psuedocode
 *
 * for each card type,
 *     x = how much would it increase my score?
 *     y = how much would it increase my opponent's score?
 *
 * choose the card type with the highest sum x+y
 *
 * getting fancy:
 * from each card type from the hand I last looked at,
 *     if my opponent chose that card type, how would that change y?
 *     come up with a probability of my opponent choosing this card type
 * use the weighted average to calculate a "smarter" y
 */

/* brute-force algo that will probably be better as well as making chopsticks decisions easier
 *
 * for every possible permutation of the choices I could make and my opponent
 * could make for the next LOOKAHEAD_LIMIT choices, which combination has the
 * highest point differential in my favor?
 *
 * this is essentially making a tree. you could kinda bfs it and cut off
 * branches that consistently are bad early, so that you have more time to spend
 * going deeper on other branches
 */

/* things I know
 * my current hand
 * my board
 * my opponent's board
 *
 * things I can remember
 * my opponent's hand
 */

type Computer struct {
	// the previous hands we have held
	// it would make sense for this to have a length equal to numPlayers - 1
	history []util.Hand
}

/* myIdx - boards[myIdx] = my board
 * boards - slice of all players' boards, not including the cards they have chosen this round
 * hand - the hand of cards I can choose from
 */
func (cp Computer) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	var originalBoards []util.Board
	copy(originalBoards, boards)

	// TODO: make an optimized scoring function that only calculates the score for a specific player
	currentScore := score.Score(boards, false)[myIdx]
	
	// potentialScores := make([]int, len(QUANTITIES))
	bestOption := -1
	highestDiff := -1
	for ct, count := range hand {
		if count < 1 {
			continue
		}
		boards[myIdx][ct]++;
		potentialDiff := score.Score(boards, roundNum == NUM_ROUNDS-1)[myIdx] - currentScore
		if potentialDiff >= highestDiff {
			bestOption = ct
			highestDiff = potentialDiff
		}
		// potentialScores[ct] = score.Score(boards, false)[myIdx]
		boards[myIdx][ct]--;
	}

	// this could cause a bug if there is some scenario in which adding a
	// card somehow decreases the total score; however, I cannot think of
	// such a scenario
	if bestOption == -1 {
		return nil, errors.New("hand has no cards")
	}

	if len(cp.history) > 0 {
		cp.history = cp.history[1:]
	}
	cp.history = append(cp.history, hand)

	// placeholder for now
	return []int{bestOption}, nil
}
