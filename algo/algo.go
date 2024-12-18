// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"log"
	. "sushigo/constants"
	"sushigo/score"
	"sushigo/util"
)

const MAX_DEPTH = 4

type (
	Computer struct {
		// what each of the boards looked like last turn
		prevBoards []util.Board

		// the previous hands we have held, most recent first
		history []util.Hand
	}

	outcome struct {
		ct        int
		depth     int
		outcomes  []*outcome
		parent    *outcome
		playerNum int // The children of this outcome are the card types that this player could choose. The parent is the card chosen before.
		scores    []int
	}
)

func getHistoryIndex(myIdx int, playerNum int, roundNum int, numPlayers int) int {
	return ((myIdx-playerNum)*(PASS_DIRECTIONS[roundNum]) + numPlayers) % numPlayers
}

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
func (cp *Computer) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	// TODO: add chopstick support

	numPlayers := len(boards)

	// update cp.history based on board changes
	// after we've seen a hand, we can know exactly what cards it contains
	for i, prevBoard := range cp.prevBoards {
		currentBoard := boards[i]
		diff := util.Hand{}
		for ct := range len(QUANTITIES) {
			diff[ct] = currentBoard.GetQuantityNoErr(ct) - prevBoard.GetQuantityNoErr(ct)
		}
		historyIndex := getHistoryIndex(myIdx, i, roundNum, numPlayers)
		if historyIndex < len(cp.history) {
			for ct, dt := range diff {
				if dt == 0 {
					continue
				}
				if util.IsNigiriOnWasabi(ct) {
					newCt, err := util.UnWasabiify(ct)
					if err != nil {
						log.Panicf("Received error while unwasabiifying: %v", err)
					} else {
						cp.history[historyIndex][newCt] -= dt
					}
				} else if ct == WASABI && dt < 0 {
					// wasabi disappeared during wasabiification;
					// we don't want to change history for that
				} else {
					cp.history[historyIndex][ct] -= dt
				}
			}
		}
	}

	cp.history = append([]util.Hand{hand}, cp.history...)
	if len(cp.history) > numPlayers {
		// the last hand in history is an older version of the first; remove it
		cp.history = cp.history[:len(cp.history)-1]
	}

	for pn := range numPlayers {
		historyIndex := getHistoryIndex(myIdx, pn, roundNum, numPlayers)
		if historyIndex < len(cp.history) {
			log.Printf("%v thinks %v has: %v\n", myIdx, pn, cp.history[historyIndex])
		}
	}

	rootOutcome := outcome{ct: -1, depth: 0, playerNum: myIdx}
	lowestScores := make(map[*outcome]int)
	next := []*outcome{&rootOutcome} // queue, first element is the next to look at
	var currentOutcome *outcome
	var preferredOutcome *outcome
	for currentOutcome == nil || (currentOutcome.depth < MAX_DEPTH && len(next) > 0) {
		// pop the front of the queue into currentChoice
		currentOutcome = next[0]
		next = next[1:]

		log.Printf("at depth %v\n", currentOutcome.depth)

		// find the hand of currentOutcome.playerNum
		var currentHand util.Hand
		historyIndex := getHistoryIndex(myIdx, currentOutcome.playerNum, roundNum, numPlayers)
		if historyIndex < len(cp.history) {
			currentHand = cp.history[historyIndex]
			// TODO: modify currentHand to remove cards played in parent outcomes
			log.Printf("%v thinks %v has %v\n", myIdx, currentOutcome.playerNum, currentHand)
		} else {
			// fill the hand with -1s to represent a hand we have not yet seen
			for i := range currentHand {
				currentHand[i] = -1
			}
		}

		// populate currentOutcomes.{outcomes, scores}
		for ct, count := range currentHand {
			// only card types with at least one card can be played
			if count == 0 || util.IsNigiriOnWasabi(ct) {
				continue
			}

			toAdd := &outcome{
				ct:        ct,
				depth:     currentOutcome.depth + 1,
				parent:    currentOutcome,
				playerNum: (currentOutcome.playerNum - 1 + numPlayers) % numPlayers,
			}


			potentialBoards := make([]util.Board, numPlayers)
			copy(potentialBoards, boards)

			// fill potentialBoards with what the boards would look
			// like if this outcome were to occur
			parent := toAdd
			for i := 0; parent != nil && parent.ct > -1; i++ {
				boardIndex := (myIdx + toAdd.depth - i + numPlayers) % numPlayers
				err := potentialBoards[boardIndex].AddCard(parent.ct)
				if err != nil {
					log.Panic(err)
				}
				parent = parent.parent
			}

			// add the scores corresponding to currentOutcome
			toAdd.scores = score.Score(potentialBoards, roundNum)
			currentOutcome.outcomes = append(currentOutcome.outcomes, toAdd)

			if toAdd.depth == MAX_DEPTH-1 {
				// find the original choice that led to this outcome
				// it will be a direct child of rootOutcomes
				searchOutcome := toAdd
				for searchOutcome.parent != &rootOutcome {
					searchOutcome = searchOutcome.parent
				}

				// is this the lowest score resulting from that outcome?
				// if so, update lowestScores
				thisScore := toAdd.scores[myIdx]
				currentMin, ok := lowestScores[searchOutcome]
				if !ok || (ok && (thisScore < currentMin)) {
					lowestScores[searchOutcome] = thisScore
				}

				// is this score the highest in lowestScores?
				// if so, this outcome is the best outcome
				highest := thisScore
				thisIsHighest := true
				for _, score := range lowestScores {
					if score > highest {
						thisIsHighest = false
						break
					}
				}
				if thisIsHighest {
					preferredOutcome = searchOutcome
				}
			}
		}

		// Add the new outcomes to next
		for _, oc := range currentOutcome.outcomes {
			next = append(next, oc)
		}
	}

	cp.prevBoards = boards

	if preferredOutcome != nil {
		return []int{preferredOutcome.ct}, nil
	} else {
		return nil, errors.New("did not find a preferred outcome")
	}
}
