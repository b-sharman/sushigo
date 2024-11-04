package main

import (
	"fmt"
	"log"
	. "sushigo/constants"
	"sushigo/player"
	"sushigo/ui"
	"sushigo/util"
)

func main() {
	num_players := ui.GetNumPlayers()
	if num_players < MIN_PLAYERS || num_players > MAX_PLAYERS {
		log.Panicf("num_players has impermissible value of %v", num_players)
	}
	cards_per_player := CARD_COUNT - num_players

	players := make([]player.Player, 0, num_players)
	// default is the first player is human and the rest are computers
	players = append(players, new(player.HumanPlayer))
	for i := 1; i < num_players; i++ {
		players = append(players, new(player.ComputerPlayer))
	}

	hands := make([]util.Hand, num_players)

	deck := NewDeck()
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
			fmt.Printf("\nPlayer %v's board:\n", j)
			util.PrintHand(util.Hand(player.GetBoard()))
		}

		for j, player := range players {
			hand_idx := ((num_players-PASS_DIRECTIONS[0])*i + j) % num_players
			ct, err := player.ChooseCard(&hands[hand_idx])
			if err != nil {
				log.Printf("Warning: the %vth player returned an error when picking a card: %v", j, err)
				continue
			}
			if util.IsNigiriOnWasabi(ct) {
				log.Printf("The %vth player tried to select a nigiri on wasabi. They should just select a nigiri instead. Their choice will be ignored.", j)
			}

			hands[hand_idx][ct]--
			if util.IsNigiri(ct) && player.GetBoard()[WASABI] > 0 {
				new_ct, err := util.Wasabiify(ct)
				if err != nil {
					log.Printf("Warning: wasabiification of ct %v (%v) failed: %v", ct, NAMES[ct], err)
				} else {
					ct = new_ct
					player.RemoveWasabi()
				}
			}
			player.AddCard(ct)
		}
	}

	boards := make([]util.Board, 0, num_players)
	for _, player := range players {
		boards = append(boards, player.GetBoard())
	}
	scores := score(boards)

	for _, score := range scores {
		fmt.Println(score)
	}
}
