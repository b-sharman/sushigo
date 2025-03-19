package plr

import (
	. "sushigo/constants"
	"sushigo/ui"
	"sushigo/util"
)

type Human struct{}

// remove a card type from the hand and add it to the board
func (hp Human) ChooseCard(roundNum int, myIdx int, boards []util.Board, hand util.Hand) ([]Card, error) {
	return ui.GetCardType(boards[myIdx].GetQuantityNoErr(CHOPSTICKS) > 0, hand), nil
}
