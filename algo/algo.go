// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"fmt"
	"log"
	_ "math"
	"slices"

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

	// a node of the search graph
	outcome struct {
		// TODO: if memory is an issue, both boards and hands could be computed from turn.cards

		// each player's board at this state
		boards []util.Board

		// a number representing how good this outcome is for us - the bigger, the better
		evaluation int

		// each player's hand at this state
		hands []util.Hand
	}

	// an edge in the search graph
	turn struct {
		// cards[playerNum] = int representing the card the player chose (will have to change to []int when chopstick support is added)
		cards []int // slice of cts

		// how many hypothetical turns it took to get to result
		depth int

		// the state of the boards before the cards were chosen
		parent *outcome

		// the state of the boards after the cards were chosen
		result *outcome
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

	hands := make([]util.Hand, 0, numPlayers)
	for pn := range numPlayers {
		historyIndex := getHistoryIndex(myIdx, pn, roundNum, numPlayers)
		if historyIndex < len(cp.history) {
			log.Printf("%v thinks %v has: %v\n", myIdx, pn, cp.history[historyIndex])
			hands = append(hands, cp.history[historyIndex])
		} else {
			// a count value of -1 means that we are unsure of whether the hand has that ct
			var placeholderHand util.Hand
			for ct := range placeholderHand {
				placeholderHand[ct] = -1
			}
			hands = append(hands, placeholderHand)
		}
	}
	log.Printf("hands: %v\n", hands)

	// queue, first element is the next to look at
	// contains the nodes of the graph described by edges
	next := []*outcome{{boards: boards, hands: hands}}
	// the paths from one element of next to another
	edges := []*turn{}
	var currentOutcome *outcome
	for depth := 1; depth <= MAX_DEPTH; depth++ {
		fmt.Println()

		// pop the front of the queue into currentOutcome
		currentOutcome = next[0]
		next = next[1:]

		// calculate all possible combinations of cards players could choose
		// TODO: make combos one slice deeper to account for chopstick use
		combos := make([][]int, 1)
		for len(combos[0]) != numPlayers {
			first := combos[0]
			combos = combos[1:]
			handIdx := ((numPlayers+PASS_DIRECTIONS[roundNum])*depth + len(first)-1) % numPlayers
			for ct, count := range currentOutcome.hands[handIdx] {
				if !util.IsNigiriOnWasabi(ct) && count != 0 {
					combos = append(combos, append(first, ct))
				}
			}
		}

		for _, choices := range combos {
			resultingBoards := make([]util.Board, 0, numPlayers)
			resultingHands := make([]util.Hand, numPlayers)
			for i, cHand := range currentOutcome.hands {
				for ct, count := range cHand {
					resultingHands[i][ct] = count
				}
			}
			for playerNum, board := range currentOutcome.boards {
				newBoard := board.DeepCopy()
				err := newBoard.AddCard(choices[playerNum])
				if err != nil {
					return nil, fmt.Errorf("error when calculating outcomes: %v", err)
				}
				resultingBoards = append(resultingBoards, newBoard)

				handIdx := ((numPlayers+PASS_DIRECTIONS[roundNum])*depth + playerNum-1) % numPlayers
				resultingHands[handIdx][choices[playerNum]]--
			}

			scores := score.Score(resultingBoards, roundNum)
			evaluation := scores[myIdx] - slices.Max(scores)

			result := outcome{
				boards: resultingBoards,
				evaluation: evaluation,
				hands: resultingHands,
			}
			// TODO: traverse the graph; if node == result, let result = node
			// consider using a hash function to accomplish this
			next = append(next, &result)

			// make an edge connecting the current outcome to the new one
			edges = append(edges, &turn{
				cards: choices,
				depth: depth,
				parent: currentOutcome,
				result: &result,
			})

			log.Printf("result: %v\n", result)
			log.Printf("added edge to depth %v\n", depth)
		}
	}

	cp.prevBoards = boards

	/*
	// for each edge
	//     if edge.depth == MAX_DEPTH
	//         if the corresponding result has the best evaluation so far
	//             make it The Node
	// find The Node's parent at depth 1
	// return the ct of the parent

	best_eval := math.MinInt
	var bestNode *outcome
	for _, edge := range edges {
		if edge.depth == MAX_DEPTH && edge.result.evaluation >= best_eval {
			bestNode = edge.result
		}
	}
	// find the ancestor of bestNode living at depth 1
	// uh oh, there's no good way to trace back our path!
	*/

	return nil, errors.New("did not find a preferred outcome")
}
