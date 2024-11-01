package main

import (
	"cmp"
	. "sushigo/constants"
	"sushigo/util"
)

func extreme_count(slc []int, comp func(a, b int) int) []int {
	ex_val := slc[0]
	ex_idc := []int{0} // indices of the most extreme values
	for i, num := range slc[1:] {
		i += 1 // correct offset introduced by slicing from 1
		if c := comp(num, ex_val); c > -1 {
			if c == 1 {
				ex_idc = nil
				ex_val = num
			}
			ex_idc = append(ex_idc, i)
		}
	}
	return ex_idc
}

// return a scores []int with the same length as boards
func score(boards []util.Board) []int {
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
	first_pudding := boards[0][PUDDING]
	all_equal := true
	puddings := make([]int, 0, len(boards))
	for _, board := range boards {
		puddings = append(puddings, board[PUDDING])
		all_equal = all_equal && board[PUDDING] == first_pudding
	}
	// no points are awarded if all players have the same number
	if !all_equal {
		// penalize players with the least puddings
		if len(boards) > 2 {
			least_p_idc := extreme_count(puddings, func(a, b int) int { return -1 * cmp.Compare(a, b) })
			for _, idx := range least_p_idc {
				scores[idx] -= 6 / len(least_p_idc)
			}
		}
		// award players with the most puddings
		most_p_idc := extreme_count(puddings, cmp.Compare)
		for _, idx := range most_p_idc {
			scores[idx] += 6 / len(most_p_idc)
		}
	}

	// makis
	makis := make([]int, 0, len(boards))
	for _, board := range boards {
		totalMakis := board[MAKI_1]*1 + board[MAKI_2]*2 + board[MAKI_3]*3
		makis = append(makis, totalMakis)
	}
	// award players with the most makis
	most_m_val := -1
	sndmost_m_val := -1
	var most_m_idc []int    // indices of the most extreme values
	var sndmost_m_idc []int // indices of the second-most extreme values
	for i, num := range makis {
		if num >= most_m_val {
			if num > most_m_val {
				// the former most is the new sndmost
				sndmost_m_val = most_m_val
				sndmost_m_idc = make([]int, len(most_m_idc))
				copy(sndmost_m_idc, most_m_idc)
				most_m_idc = nil
				most_m_val = num
			}
			most_m_idc = append(most_m_idc, i)
		} else if num >= sndmost_m_val {
			if num > sndmost_m_val {
				sndmost_m_val = num
				sndmost_m_idc = nil
			}
			sndmost_m_idc = append(sndmost_m_idc, i)
		}
	}
	if most_m_val > 0 {
		for _, idx := range most_m_idc {
			scores[idx] += 6 / len(most_m_idc)
		}
	}
	if len(most_m_idc) < 2 && sndmost_m_val > 0 {
		for _, idx := range sndmost_m_idc {
			scores[idx] += 3 / len(sndmost_m_idc)
		}
	}

	return scores
}
