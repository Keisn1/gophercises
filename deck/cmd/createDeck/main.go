package main

import (
	"deck"
	"fmt"
)

func main() {
	opts := deck.NewOpts().
		WithSort(deck.AbsRank).
		WithReverse().
		WithJokers(4).
		WithFilterRanks([]deck.Rank{deck.Ace}).
		WithFilterSuits([]deck.Suit{deck.Diamond}).
		WithMultipleDecks(4).
		DoShuffle()
}
