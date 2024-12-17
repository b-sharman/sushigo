// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"fmt"
	"log"
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
		prevBoards []util.Board
		// the previous hands we have held, most recent first
		// it would make sense for this to have a length equal to numPlayers - 1
		history []util.Hand
	}

	Outcome struct {
		ct        int
		depth     int
		outcomes  []*Outcome
		parent    *Outcome
		playerNum int // The children of this outcome are the card types that this player could choose. The parent is the card chosen before.
		scores    []int
	}
)

/* myIdx - boards[myIdx] = my board
 * boards - slice of all players' boards, not including the cards they have chosen this round
 * hand - the hand of cards I can choose from
 */
func (cp *Computer) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	// TODO: add chopstick support

	numPlayers := len(boards)

	fmt.Println(myIdx)

	// remove whatever new cards are on the boards from cp.history
	// after we've seen a hand, we can know exactly what cards it contains
	for i, prevBoard := range cp.prevBoards {
		currentBoard := boards[i]
		diff := util.Board{}
		fmt.Printf("prevBoard: %v\n", prevBoard)
		fmt.Printf("currentBoard: %v\n", currentBoard)
		for ct := range currentBoard {
			diff[ct] = currentBoard[ct] - prevBoard[ct]
		}
		historyIndex := ((myIdx - i)*(PASS_DIRECTIONS[roundNum]) + numPlayers) % numPlayers
		fmt.Printf("board %v, diff %v, histIdx %v\n", i, diff, historyIndex)
		if historyIndex < len(cp.history) {
			fmt.Printf("before: %v\n", cp.history[historyIndex])
			for ct, dt := range diff {
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
					fmt.Printf("skipped wasabi diff (%v)\n", dt)
				} else {
					cp.history[historyIndex][ct] -= dt
				}
			}
			fmt.Printf("after (%v): %v\n", i, cp.history[historyIndex])
		}
	}

	cp.history = append([]util.Hand{hand}, cp.history...)
	if len(cp.history) > numPlayers {
		// the last hand in history is an older version of the first; remove it
		cp.history = cp.history[:len(cp.history)-1]
	}

	rootOutcome := Outcome{ct: -1, depth: 0, playerNum: myIdx}
	lowestScores := make(map[*Outcome]int)
	next := []*Outcome{&rootOutcome} // queue, first element is the next to look at
	var currentOutcome *Outcome
	var preferredOutcome *Outcome
	for currentOutcome == nil || (currentOutcome.depth < MAX_DEPTH && len(next) > 0) {
		// pop the front of the queue into currentChoice and push its children to the back
		currentOutcome = next[0]
		next = next[1:]

		// find the hand of this player
		var currentHand util.Hand
		// fill the hand with -1s by default to represent a hand we have not yet seen
		for i := range currentHand {
			currentHand[i] = -1
		}
		historyIndex := ((myIdx - currentOutcome.playerNum)*(PASS_DIRECTIONS[roundNum]) + numPlayers) % numPlayers
		if historyIndex < len(cp.history) {
			currentHand = cp.history[historyIndex]
			fmt.Printf("%v thinks %v has: %v\n", myIdx, currentOutcome.playerNum, currentHand)
		}

		// populate currentOutcomes.{outcomes, scores}
		for ct, count := range currentHand {
			// only card types with at least one card can be played
			if count == 0 {
				continue
			}

			// TODO: wasabiify ct if potentialBoards[myIdx] has a wasabi

			toAdd := &Outcome{
				ct:        ct,
				depth:     currentOutcome.depth + 1,
				parent:    currentOutcome,
				playerNum: (currentOutcome.playerNum - 1 + numPlayers) % numPlayers,
			}

			// add the scores corresponding to currentOutcome
			potentialBoards := make([]util.Board, numPlayers)
			copy(potentialBoards, boards)
			// add the ct of this outcome and all its parents to the boards
			potentialBoards[(myIdx+toAdd.depth)%numPlayers][ct]++
			parent := toAdd.parent
			for i := 0; parent != nil && parent.ct > -1; i++ {
				boardIndex := (myIdx + toAdd.depth - i + numPlayers) % numPlayers
				potentialBoards[boardIndex][parent.ct]++
				parent = parent.parent
			}
			toAdd.scores = score.Score(potentialBoards, roundNum == NUM_ROUNDS-1)
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
