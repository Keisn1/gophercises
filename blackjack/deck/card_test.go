package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	c := Card{Rank: Ace, Suit: Diamond}
	fmt.Println(c.String())

	// Output:
	// Ace of Diamonds
}

func TestLen(t *testing.T) {
	want := 52
	got := len(New())
	if got != want {
		t.Errorf("len(New()) = \"%v\"; want \"%v\"", got, want)
	}
}

func TestDefaultSort(t *testing.T) {
	want := Card{Rank: Ace, Suit: Diamond}
	got := New()[0]
	if got != want {
		t.Errorf("New(WithSort(absRank))[0] = \"%v\"; want \"%v\"", got.String(), want.String())
	}
}

func TestReverseSort(t *testing.T) {
	opts := NewOpts().WithReverse()
	got := NewWithOpts(opts)[0]
	want := Card{Rank: King, Suit: Club}
	if got != want {
		t.Errorf("New(WithReverse())[0] = \"%v\"; want \"%v\"", got.String(), want.String())
	}
}

func TestNbrOfJokers(t *testing.T) {
	want := 10
	opts := NewOpts().WithJokers(10)
	deck := NewWithOpts(opts)
	count := 0
	for _, c := range deck {
		if c.Suit == Joker {
			count++
		}
	}
	if count != want {
		t.Errorf("Number of Jokers in Deck = \"%v\"; want \"%v\"", count, want)
	}
}

func TestFilterRank(t *testing.T) {
	opts := NewOpts().WithFilterRanks([]Rank{Two, Five})
	deck := NewWithOpts(opts)
	for _, c := range deck {
		if c.Rank == Two || c.Rank == Five {
			t.Errorf("Found %s in Deck, don't want %v", c.String(), []Rank{Two, Five})
		}
	}
}

func TestFilterSuit(t *testing.T) {
	opts := NewOpts().WithFilterSuits([]Suit{Diamond})
	deck := NewWithOpts(opts)
	for _, c := range deck {
		if c.Suit == Diamond {
			t.Errorf("Found %s in Deck, don't want %v", c.String(), []Suit{Diamond})
		}
	}
}

func TestMultipleDecks(t *testing.T) {
	want := 4 * 52
	opts := NewOpts().WithMultipleDecks(4)
	got := len(NewWithOpts(opts))
	if got != want {
		t.Errorf("len(New1(MultipleDecks(4))) = \"%v\"; want \"%v\"", got, want)
	}
}

func TestShuffle(t *testing.T) {
	opts := NewOpts().DoShuffle()
	got := NewWithOpts(opts)
}
