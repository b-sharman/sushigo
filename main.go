package main

import (
	"fmt"
	"log"
	. "sushigo/constants"
	"sushigo/plr"
	"sushigo/ui"
	"sushigo/util"
)

func playRound(deck *Deck, players []*plr.Player, cardsPerPlayer int) []util.Board {
	numPlayers := len(players)

	for _, player := range players {
		plr.ClearBoard(player)
	}

	hands := make([]util.Hand, numPlayers)
	// deal as many hands as there are players
	for i := range hands {
		cards, err := deck.NextNCards(cardsPerPlayer)
		if err != nil {
			log.Panic(err)
		}
		for _, ct := range cards {
			hands[i][ct]++
		}
	}

	// let players pick cards until the hands are exhausted
	for i := 0; i < cardsPerPlayer; i++ {
		for j, player := range players {
			handIdx := ((numPlayers-PASS_DIRECTIONS[0])*i + j) % numPlayers
			cts, err := plr.ChooseCard(player, &hands[handIdx])
			if err != nil {
				log.Printf("Warning: the %vth player returned an error when picking a card: %v", j, err)
				continue
			}

			if len(cts) > 1 {
				// Chopsticks used
				err := plr.RemoveCard(player, CHOPSTICKS)
				if err != nil {
					log.Printf("Player %v tried to play two cards but has no chopsticks. Only the first card will be considered.", j)
					cts = cts[:1]
				} else {
					// add the player's chopsticks back into the hand
					hands[handIdx][CHOPSTICKS]++
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

				hands[handIdx][ct]--
				if util.IsNigiri(ct) && plr.GetBoard(player)[WASABI] > 0 {
					newCt, err := util.Wasabiify(ct)
					if err != nil {
						log.Printf("Warning: wasabiification of ct %v (%v) failed: %v", ct, NAMES[ct], err)
					} else {
						ct = newCt
						plr.RemoveCard(player, WASABI)
					}
				}
				plr.AddCard(player, ct)
			}
		}

		fmt.Println()
		for j, player := range players {
			fmt.Printf("Player %v's board:\n", j)
			util.PrintHand(util.Hand(plr.GetBoard(player)))
		}
	}

	boards := make([]util.Board, 0, numPlayers)
	for _, player := range players {
		boards = append(boards, plr.GetBoard(player))
	}

	return boards
}

func main() {
	numPlayers := ui.GetNumPlayers()
	if numPlayers < MIN_PLAYERS || numPlayers > MAX_PLAYERS {
		log.Panicf("numPlayers has impermissible value of %v", numPlayers)
	}
	cardsPerPlayer := CARD_COUNT - numPlayers

	players := make([]*plr.Player, 0, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players = append(players, new(plr.Player))
	}
	// make the first player human
	players[0].IsHuman = true

	var scores [][]int

	deck := NewDeck()

	for i := 0; i < NUM_ROUNDS; i++ {
		boards := playRound(&deck, players, cardsPerPlayer)
		roundScores := score(boards, i == NUM_ROUNDS-1)
		scores = append(scores, roundScores)

		ui.PrintScores(scores, numPlayers, i)
	}
}
