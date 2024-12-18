package main

import (
	"fmt"
	"log"
	"sushigo/algo"
	. "sushigo/constants"
	"sushigo/plr"
	"sushigo/score"
	"sushigo/ui"
	"sushigo/util"
)

func playRound(roundNum int, deck *Deck, players []*plr.Player, cardsPerPlayer int) []util.Board {
	numPlayers := len(players)

	for _, player := range players {
		player.Board.Clear()
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
		// selected cards must be stored in a queue so that computer
		// players can't see boards ahead of time
		//
		// no need to implement a real queue, we can just use a slice
		addQueue := make([][]int, 0, numPlayers)
		for j, player := range players {
			handIdx := ((numPlayers+PASS_DIRECTIONS[roundNum])*i + j) % numPlayers
			fmt.Printf("main: player %v has hand: %v\n", j, hands[handIdx])
			cts, err := player.Chooser.ChooseCard(roundNum, j, plr.BoardsFromPlayers(players), hands[handIdx])
			if err != nil {
				log.Printf("Warning: the %vth player returned an error when picking a card: %v", j, err)
				continue
			}
			for _, ct := range cts {
				if ct < 0 || ct >= len(QUANTITIES) {
					log.Printf("Warning: the %vth player returned invalid card type %v", j, ct)
					cts = nil
				}
				if hands[handIdx][ct] < 1 {
					log.Printf("Warning: the %vth player requested card type %v (%v), but there are no such cards in the hand", j, ct, NAMES[ct])
					cts = nil
				}
			}
			addQueue = append(addQueue, cts)

			if len(cts) > 1 {
				// Chopsticks used
				err := player.Board.RemoveCard(CHOPSTICKS)
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
				hands[handIdx][ct]--
			}
		}

		for j, player := range players {
			cts := addQueue[j]
			for _, ct := range cts {
				err := player.Board.AddCard(ct)
				if err != nil {
					log.Printf("Warning: failed to add ct %v to player %v: %v", ct, j, err)
				}
			}
		}

		fmt.Println()
		for j, player := range players {
			fmt.Printf("Player %v's board:\n", j)
			util.PrintHand(player.Board.ToHand())
		}
	}

	return plr.BoardsFromPlayers(players)
}

func main() {
	numPlayers := ui.GetNumPlayers()
	if numPlayers < MIN_PLAYERS || numPlayers > MAX_PLAYERS {
		log.Panicf("numPlayers has impermissible value of %v", numPlayers)
	}
	cardsPerPlayer := CARD_COUNT - numPlayers

	players := make([]*plr.Player, 0, numPlayers)
	for i := 0; i < numPlayers; i++ {
		newPlayer := new(plr.Player)
		if i == 0 {
			newPlayer.Chooser = &plr.Human{}
		} else {
			newPlayer.Chooser = &algo.Computer{}
		}
		players = append(players, newPlayer)
	}

	var scores [][]int

	deck := NewDeck()

	for i := 0; i < NUM_ROUNDS; i++ {
		boards := playRound(i, &deck, players, cardsPerPlayer)
		roundScores := score.Score(boards, i == NUM_ROUNDS-1)
		scores = append(scores, roundScores)

		ui.PrintScores(scores, numPlayers, i)
	}
}
