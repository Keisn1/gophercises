// package deckOld

// import (
// 	"fmt"
// )

// // My initial solution
// type By func(c1, c2 *Card) bool

// func (by By) Sort(deck []Card) {
// 	ds := &DeckSorter{
// 		deck:   deck,
// 		sortBy: by,
// 	}
// 	sort.Sort(ds)
// }

// type DeckSorter struct {
// 	deck   []Card
// 	sortBy func(c1, c2 *Card) bool
// }

// func (s *DeckSorter) Len() int {
// 	return len(s.deck)
// }

// func (s *DeckSorter) Swap(i, j int) {
// 	s.deck[i], s.deck[j] = s.deck[j], s.deck[i]
// }

// func (s *DeckSorter) Less(i, j int) bool {
// 	return s.sortBy(&s.deck[i], &s.deck[j])
// }
// func New(opts ...func(*DeckSorter)) []Card {
// 	var deck []Card
// 	for _, rank := range ranks {
// 		for _, suit := range suits {
// 			deck = append(deck, Card{
// 				Suit: suit,
// 				Rank: rank,
// 			})
// 		}
// 	}

// 	ds := DeckSorter{deck: deck}
// 	for _, opt := range opts {
// 		opt(&ds)
// 	}

// 	By(ds.sortBy).Sort(deck)
// 	return deck
// }

// func DefaultComp(c1, c2 *Card) bool {
// 	if c1.Suit < c2.Suit {
// 		return true
// 	} else if c1.Suit > c2.Suit {
// 		return false
// 	}
// 	return c1.Rank < c2.Rank
// }

// func WithSorting(comp func(c1, c2 *Card) bool) func(*DeckSorter) {
// 	return func(ds *DeckSorter) {
// 		ds.sortBy = comp
// 	}
// }
