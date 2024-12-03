// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"fmt"
	. "sushigo/constants"
	"sushigo/score"
	"sushigo/util"
)

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

// TODO: move to constants
const MAX_DEPTH = 4

type (
	Computer struct {
		// the previous hands we have held
		// it would make sense for this to have a length equal to numPlayers - 1
		history []util.Hand
	}

	Outcome struct {
		ct        int
		outcomes  []*Outcome
		parent    *Outcome
		playerNum int // The children of this outcome are the card types that this player could choose. The parent is the card chosen before.
		scores    []int
	}
)

var EVERYTHINGHAND = util.Hand{
	1, // CHOPSTICKS
	1, // DUMPLING
	1, // MAKI_1
	1, // MAKI_2
	1, // MAKI_3
	1, // NIGIRI_1
	1, // NIGIRI_2
	1, // NIGIRI_3
	0, // NIGIRI_1_ON_WASABI
	0, // NIGIRI_2_ON_WASABI
	0, // NIGIRI_3_ON_WASABI
	1, // PUDDING
	1, // SASHIMI
	1, // TEMPURA
	1, // WASABI
}

/* myIdx - boards[myIdx] = my board
 * boards - slice of all players' boards, not including the cards they have chosen this round
 * hand - the hand of cards I can choose from
 */
func (cp *Computer) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	// TODO: add chopstick support

	numPlayers := len(boards)

	// TODO: keep track of the history of boards and use the differences to
	// modify the cp.history to account for the cards that were removed
	// from hands we saw previously
	cp.history = append(cp.history, hand)
	if len(cp.history) > numPlayers {
		// the first element of history is an older version of the last; remove it
		cp.history = cp.history[1:]
	}

	rootOutcome := Outcome{ct: -1, playerNum: myIdx}
	// // TODO: do we have to default-initialize rootOutcome.outcomes with playerNum: len(myIdx+1) % numPlayers ?
	// for ct, count := range hand {
	// 	if count > 0 {
	// 		rootOutcome.outcomes = append(rootOutcome.outcomes, &Outcome{ct: ct})
	// 	}
	// }

	next := []*Outcome{&rootOutcome} // queue, first element is the next to look at
	var currentOutcome *Outcome
	// points to the edge node outcome with the highest minimum score
	var highestLowest *Outcome
	highestScore := -1
	for depth := 0; depth < MAX_DEPTH && len(next) > 0; depth++ {
		// pop the front of the queue into currentChoice and push its children to the back
		currentOutcome = next[0]
		next = next[1:]

		// find the hand of this player
		// Currently we assume it has all the cards it had when we last saw it.
		currentHand := EVERYTHINGHAND
		historyIndex := (currentOutcome.playerNum - myIdx + numPlayers) % numPlayers
		if historyIndex < len(cp.history) {
			currentHand = cp.history[historyIndex]
		}

		// populate currentOutcomes.{outcomes, scores}
		for ct, count := range currentHand {
			// only card types with at least one card can be played
			if count < 1 {
				continue
			}

			maxScore := -1
			toAdd := &Outcome{
				ct:        ct,
				parent:    currentOutcome,
				playerNum: (currentOutcome.playerNum - 1 + numPlayers) % numPlayers,
			}

			// add the scores corresponding to currentOutcome
			potentialBoards := make([]util.Board, numPlayers)
			copy(potentialBoards, boards)
			// add the ct of this outcome and all its parents to the boards
			potentialBoards[(myIdx+depth)%numPlayers][ct]++
			parent := toAdd.parent
			for i := 0; parent != nil && parent.ct > -1; i++ {
				potentialBoards[(myIdx+depth-i+numPlayers)%numPlayers][parent.ct]++
				parent = parent.parent
			}
			toAdd.scores = score.Score(potentialBoards, roundNum == NUM_ROUNDS-1)
			currentOutcome.outcomes = append(currentOutcome.outcomes, toAdd)

			if toAdd.scores[myIdx] > maxScore {
				maxScore = toAdd.scores[myIdx]
			}
			if maxScore > highestScore {
				highestScore = maxScore
				minScore := -1
				var minOutcome *Outcome
				for _, oc := range currentOutcome.outcomes {
					if minScore == -1 || oc.scores[myIdx] < minScore {
						minScore = oc.scores[myIdx]
						minOutcome = oc
					}
				}
				highestLowest = minOutcome
			}
		}

		// Add the new outcomes to next
		for _, oc := range currentOutcome.outcomes {
			next = append(next, oc)
		}
	}

	fmt.Print(*highestLowest)
	parent := highestLowest.parent
	for {
		if parent.parent == nil {
			break
		}
		parent = parent.parent
	}

	return []int{parent.ct}, nil
}
