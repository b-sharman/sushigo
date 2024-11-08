package plr

import (
	. "sushigo/constants"
	"sushigo/ui"
	"sushigo/util"
)

type Human struct{}

// remove a card type from the hand and add it to the board
func (hp Human) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]int, error) {
	hasChopsticks := boards[myIdx][CHOPSTICKS] > 0
	return ui.GetCardType(hasChopsticks, hand), nil
}
