// An algorithmic chooser that tries to always choose the best card.

package algo

import (
	"errors"
	"fmt"
	"log"
	"math"
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
		from *outcome

		// the state of the boards after the cards were chosen
		to *outcome
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

	// the nodes of the graph described by edges
	nodes := []*outcome{{boards: boards, hands: hands}}
	nextIdx := 0
	// the paths from one node to another
	edges := make(map[*outcome][]*turn)
	var currentOutcome *outcome
	bestEval := math.MinInt
	var bestNode *outcome
	for depth := 1; depth <= MAX_DEPTH; depth++ {
		fmt.Println()

		// pop the front of the queue into currentOutcome
		currentOutcome = nodes[nextIdx]

		// calculate all possible combinations of cards players could choose

		// There are len(combos) different possibilities for the cards
		// players could choose this turn. combos[i] has length
		// numPlayers and represents the cts of each player for that
		// possibility.

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
			if len(combos) < 1 {
				break
			}
		}

		for _, choices := range combos {
			resultingBoards := make([]util.Board, 0, numPlayers)
			resultingHands := make([]util.Hand, numPlayers)
			// copy currentOutcome.hands into resultingHands
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
			nodes = append(nodes, &result)
			nextIdx++

			if depth == MAX_DEPTH && evaluation > bestEval {
				bestNode = &result
			}

			// make an edge connecting the current outcome to the new one
			edge := &turn{
				cards: choices,
				depth: depth,
				from: currentOutcome,
				to: &result,
			}
			edges[currentOutcome] = append(edges[currentOutcome], edge)
			edges[&result] = append(edges[&result], edge)

			log.Printf("result: %v\n", result)
			log.Printf("added edge to depth %v\n", depth)
		}
	}

	cp.prevBoards = boards

	// find the ancestor of bestNode living at depth 1
	// ancestors with depth != 1
	interimAncestors := []*outcome{bestNode}
	// ancestors with depth == 1
	var finalAncestors []*outcome
	for len(interimAncestors) > 0 {
		current := interimAncestors[0]
		interimAncestors = interimAncestors[1:]
		correspondingEdges, ok := edges[current]
		if ok {
			for _, edge := range correspondingEdges {
				// if it took 2 hypothetical turns to get to `result`, it took 1 hypothetical turn to get to `from`
				if edge.depth == 2 {
					finalAncestors = append(finalAncestors, edge.from)
				} else {
					interimAncestors = append(interimAncestors, edge.from)
				}
			}
		}
	}
	for i, fa := range finalAncestors {
		log.Printf("finalAncestors[%v] = %v\n", i, fa)
	}

	if len(finalAncestors) < 1 {
		return nil, errors.New("did not find a preferred outcome")
	}
	return []int{edges[finalAncestors[0]][0].cards[myIdx]}, nil
}
