package main

import (
	"fmt"
	"log"
	. "sushigo/constants"
	plr "sushigo/player"
	"sushigo/ui"
	"sushigo/util"
)

func playRound(deck *Deck, players []*plr.Player, cards_per_player int) []util.Board {
	num_players := len(players)

	for _, player := range players {
		plr.ClearBoard(player)
	}

	hands := make([]util.Hand, num_players)
	// deal as many hands as there are players
	for i := range hands {
		cards, err := deck.NextNCards(cards_per_player)
		if err != nil {
			log.Panic(err)
		}
		for _, ct := range cards {
			hands[i][ct]++
		}
	}

	// let players pick cards until the hands are exhausted
	for i := 0; i < cards_per_player; i++ {
		for j, player := range players {
			hand_idx := ((num_players-PASS_DIRECTIONS[0])*i + j) % num_players
			cts, err := plr.ChooseCard(player, &hands[hand_idx])
			if err != nil {
				log.Printf("Warning: the %vth player returned an error when picking a card: %v", j, err)
				continue
			}

			if len(cts) > 1 {
				// Chopsticks used
				err := plr.RemoveChopsticks(player)
				if err != nil {
					log.Printf("Player %v tried to play two cards but has no chopsticks. Only the first card will be considered.", j)
					cts = cts[:1]
				} else {
					// add the player's chopsticks back into the hand
					hands[hand_idx][CHOPSTICKS]++
				}

				// At first I thought, if wasabi and nigiri are chosen
				// simultaneously with chopsticks, does that mean the nigiri is
				// forced to be on the wasabi, as usual? The answer is no; the
				// player can choose which to play first. If they already have an
				// uncovered wasabi, the order doesn't matter. Otherwise, they can
				// play the nigiri first in order to not apply it to the wasabi, or
				// they can play the wasabi first, forcing themselves to wasabiify
				// their nigiri.
				//
				// Phil Walker-Harding, the creator of Sushi Go, verified this on
				// Board Game Geek.
				// https://boardgamegeek.com/thread/1014756/chopsticks-with-wasabi-and-nigiri
				//
				// That's nice for me; it means there's one less thing to implement.
			}

			for _, ct := range cts {
				if util.IsNigiriOnWasabi(ct) {
					log.Printf("The %vth player tried to select a nigiri on wasabi. They should just select a nigiri instead. Their choice will be ignored.", j)
					continue
				}

				hands[hand_idx][ct]--
				if util.IsNigiri(ct) && plr.GetBoard(player)[WASABI] > 0 {
					new_ct, err := util.Wasabiify(ct)
					if err != nil {
						log.Printf("Warning: wasabiification of ct %v (%v) failed: %v", ct, NAMES[ct], err)
					} else {
						ct = new_ct
						plr.RemoveWasabi(player)
					}
				}
				plr.AddCard(player, ct)
			}
		}

		for j, player := range players {
			fmt.Printf("\nPlayer %v's board:\n", j)
			util.PrintHand(util.Hand(plr.GetBoard(player)))
		}
	}

	boards := make([]util.Board, 0, num_players)
	for _, player := range players {
		boards = append(boards, plr.GetBoard(player))
	}

	return boards
}

func main() {
	num_players := ui.GetNumPlayers()
	if num_players < MIN_PLAYERS || num_players > MAX_PLAYERS {
		log.Panicf("num_players has impermissible value of %v", num_players)
	}
	cards_per_player := CARD_COUNT - num_players

	players := make([]*plr.Player, 0, num_players)
	for i := 0; i < num_players; i++ {
		players = append(players, new(plr.Player))
	}
	// make the first player human
	players[0].IsHuman = true

	var scores [][]int

	deck := NewDeck()

	for i := 0; i < NUM_ROUNDS; i++ {
		boards := playRound(&deck, players, cards_per_player)
		roundScores := score(boards, i == NUM_ROUNDS-1)
		scores = append(scores, roundScores)

		ui.PrintScores(scores, num_players, i)
	}
}
