// An algorithmic chooser that tries to always choose the best card.

package algo

import (
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
func (cp Computer) ChooseCard(myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	var originalBoards []util.Board
	copy(originalBoards, boards)

	// software engineering problem - we need to remember what our last hand looked like; how do we do that?
	// does algo have to become a struct? can we add a "history" filed to plr.Player?

	// placeholder for now
	return []int{}, nil
}
