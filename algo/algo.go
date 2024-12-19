// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"log"
	. "sushigo/constants"
	"sushigo/score"
	"sushigo/util"
)

const MAX_DEPTH = 2

type (
	Computer struct {
		// what each of the boards looked like last turn
		prevBoards []util.Board

		// the previous hands we have held, most recent first
		history []util.Hand
	}

	outcome struct {
		// playerNum chose this card type
		ct        int

		// how many hypothetical moves have been made so far. 0 indicates the most recent move that actually happened
		depth     int

		// the outcomes this outcome leads to
		outcomes  []*outcome

		// the outcome that led to this one
		parent    *outcome

		// the player that played the card in this outcome
		playerNum int

		// the scores resulting from this outcome - empty for all but the leaf nodes
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

	lowestScores := make(map[*outcome]int)
	// queue, first element is the next to look at
	next := []*outcome{{ct: -1, depth: 0, playerNum: -1}}
	var currentOutcome *outcome
	var preferredOutcome *outcome
	var prevParent *outcome
	for currentOutcome == nil || (len(next) > 0 && next[0].depth <= MAX_DEPTH) {
		// pop the front of the queue into currentOutcome
		currentOutcome = next[0]
		next = next[1:]

		if currentOutcome != nil && currentOutcome.parent != prevParent {
			log.Println("looking at child(ren) of this outcome:")
			log.Printf("depth: %v\n", next[0].parent.depth)
			log.Printf("playerNum: %v\n", next[0].parent.playerNum)
			log.Printf("ct: %v (%v)\n", next[0].parent.ct, NAMES[next[0].parent.ct])
		}

		// set to hand when no outcomes have been calculated;
		// overridden to be the generated hand at a hypothetical
		// outcome otherwise
		currentHand := hand

		if currentOutcome.depth > 0 {
			// find the hand of currentOutcome.playerNum
			historyIndex := getHistoryIndex(
				myIdx,
				(currentOutcome.playerNum-(currentOutcome.depth-1)*PASS_DIRECTIONS[roundNum]+numPlayers)%numPlayers,
				roundNum,
				numPlayers,
			)
			if historyIndex < len(cp.history) {
				currentHand = cp.history[historyIndex]
				// modify currentHand to remove cards played in parent outcomes
				parent := currentOutcome.parent
				for i := 0; parent != nil && parent.ct > -1; i++ {
					if parent.playerNum == currentOutcome.playerNum - PASS_DIRECTIONS[roundNum] {
						currentHand[parent.ct]--
					}
					parent = parent.parent
				}
				log.Printf("%v thinks %v has %v\n", myIdx, currentOutcome.playerNum, currentHand)
			}

			log.Printf(
				"at depth %v, player %v takes %v (%v)\n",
				currentOutcome.depth,
				currentOutcome.playerNum,
				currentOutcome.ct,
				NAMES[currentOutcome.ct],
			)
		}

		// populate currentOutcomes.{outcomes, scores}
		for ct, count := range currentHand {
			// only card types with at least one card can be played
			if count == 0 || util.IsNigiriOnWasabi(ct) {
				continue
			}

			var toAdd *outcome
			if currentOutcome.depth == 0 {
				toAdd = &outcome{
					ct:        ct,
					depth:     1,
					parent:    nil,
					playerNum: myIdx,
				}
			} else {
				taDepth := currentOutcome.depth
				if currentOutcome.playerNum == (myIdx+PASS_DIRECTIONS[roundNum]+numPlayers)%numPlayers {
					taDepth++
				}
				toAdd = &outcome{
					ct:        ct,
					depth:     taDepth,
					parent:    currentOutcome,
					playerNum: (currentOutcome.playerNum - PASS_DIRECTIONS[roundNum] + numPlayers) % numPlayers,
				}
			}

			if toAdd.depth == MAX_DEPTH+1 {
				// calculate what the boards would look like if
				// this outcome were to occur
				potentialBoards := make([]util.Board, numPlayers)
				copy(potentialBoards, boards)
				parent := toAdd
				for i := 0; parent != nil && parent.ct > -1; i++ {
					err := potentialBoards[parent.playerNum].AddCard(parent.ct)
					if err != nil {
						log.Panic(err)
					}
					parent = parent.parent
				}

				// add the scores corresponding to currentOutcome
				toAdd.scores = score.Score(potentialBoards, roundNum)

				// find the original choice that led to this outcome
				searchOutcome := toAdd
				for searchOutcome.parent.depth != 1 {
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

			currentOutcome.outcomes = append(currentOutcome.outcomes, toAdd)
		}

		// Add the new outcomes to next
		for _, oc := range currentOutcome.outcomes {
			next = append(next, oc)
		}
		prevParent = currentOutcome.parent
	}

	cp.prevBoards = boards

	if preferredOutcome != nil {
		return []int{preferredOutcome.ct}, nil
	} else {
		return nil, errors.New("did not find a preferred outcome")
	}
}
