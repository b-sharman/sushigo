package main

import (
	"cmp"
	. "sushigo/constants"
	"sushigo/util"
)

func extremeCount(slc []int, comp func(a, b int) int) []int {
	exVal := slc[0]
	exIdc := []int{0} // indices of the most extreme values
	for i, num := range slc[1:] {
		i += 1 // correct offset introduced by slicing from 1
		if c := comp(num, exVal); c > -1 {
			if c == 1 {
				exIdc = nil
				exVal = num
			}
			exIdc = append(exIdc, i)
		}
	}
	return exIdc
}

// return a scores []int with the same length as boards
func score(boards []util.Board, scorePuddings bool) []int {
	scores := make([]int, len(boards))

	// types of cards that don't depend on other players
	for i, board := range boards {
		scores[i] += []int{0, 1, 3, 6, 10, 15}[board[DUMPLING]]

		scores[i] += board[NIGIRI_1] * 1
		scores[i] += board[NIGIRI_2] * 2
		scores[i] += board[NIGIRI_3] * 3
		scores[i] += board[NIGIRI_1_ON_WASABI] * 3
		scores[i] += board[NIGIRI_2_ON_WASABI] * 6
		scores[i] += board[NIGIRI_3_ON_WASABI] * 9

		scores[i] += 10 * (board[SASHIMI] / 3)
		scores[i] += 5 * (board[TEMPURA] / 2)
	}

	// types of cards that depend on other players

	// puddings
	if scorePuddings {
		firstPudding := boards[0][PUDDING]
		allEqual := true
		puddings := make([]int, 0, len(boards))
		for _, board := range boards {
			puddings = append(puddings, board[PUDDING])
			allEqual = allEqual && board[PUDDING] == firstPudding
		}
		// no points are awarded if all players have the same number
		if !allEqual {
			// penalize players with the least puddings
			if len(boards) > 2 {
				leastP_idc := extremeCount(puddings, func(a, b int) int { return -1 * cmp.Compare(a, b) })
				for _, idx := range leastP_idc {
					scores[idx] -= 6 / len(leastP_idc)
				}
			}
			// award players with the most puddings
			mostP_idc := extremeCount(puddings, cmp.Compare)
			for _, idx := range mostP_idc {
				scores[idx] += 6 / len(mostP_idc)
			}
		}
	}

	// makis
	makis := make([]int, 0, len(boards))
	for _, board := range boards {
		totalMakis := board[MAKI_1]*1 + board[MAKI_2]*2 + board[MAKI_3]*3
		makis = append(makis, totalMakis)
	}
	// award players with the most makis
	mostVal := -1
	sndmostVal := -1
	var mostIdc []int    // indices of the most extreme values
	var sndmostIdc []int // indices of the second-most extreme values
	for i, num := range makis {
		if num >= mostVal {
			if num > mostVal {
				// the former most is the new sndmost
				sndmostVal = mostVal
				sndmostIdc = make([]int, len(mostIdc))
				copy(sndmostIdc, mostIdc)
				mostIdc = nil
				mostVal = num
			}
			mostIdc = append(mostIdc, i)
		} else if num >= sndmostVal {
			if num > sndmostVal {
				sndmostVal = num
				sndmostIdc = nil
			}
			sndmostIdc = append(sndmostIdc, i)
		}
	}
	if mostVal > 0 {
		for _, idx := range mostIdc {
			scores[idx] += 6 / len(mostIdc)
		}
	}
	if len(mostIdc) < 2 && sndmostVal > 0 {
		for _, idx := range sndmostIdc {
			scores[idx] += 3 / len(sndmostIdc)
		}
	}

	return scores
}
