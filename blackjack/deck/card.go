//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	// "time"
)

type Suit uint8

const (
	Diamond Suit = iota
	Heart
	Spade
	Club
	Joker
)

var suits = [...]Suit{Diamond, Heart, Spade, Club}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var ranks = [...]Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

type Opts struct {
	rankFuncForSort func(c Card) int
	Reverse         bool
	NbrOfJokers     int
	NbrOfDecks      int
	SuitsToFilter   []Suit
	RanksToFilter   []Rank
	WithShuffle     bool
}

// Returns Default Opts, to be used with NewWithOpts
func NewOpts() Opts {
	return Opts{
		rankFuncForSort: AbsRank,
		Reverse:         false,
		NbrOfJokers:     0,
		NbrOfDecks:      1,
		WithShuffle:     false,
		SuitsToFilter:   []Suit{},
		RanksToFilter:   []Rank{},
	}
}

func (o Opts) WithSort(rankFuncForSort func(c Card) int) Opts {
	o.rankFuncForSort = rankFuncForSort
	return o
}

func (o Opts) WithReverse() Opts {
	o.Reverse = true
	return o
}

func (o Opts) WithJokers(nbrOfJokers int) Opts {
	o.NbrOfJokers = nbrOfJokers
	return o
}

func (o Opts) WithNbrOfDecks(nbrOfDecks int) Opts {
	o.NbrOfDecks = nbrOfDecks
	return o
}

func (o Opts) WithFilterSuits(suitsToFilter []Suit) Opts {
	o.SuitsToFilter = suitsToFilter
	return o
}

func (o Opts) WithFilterRanks(ranksToFilter []Rank) Opts {
	o.RanksToFilter = ranksToFilter
	return o
}

func (o Opts) DoShuffle() Opts {
	o.WithShuffle = true
	return o
}

func (o Opts) WithMultipleDecks(nbrOfDecks int) Opts {
	o.NbrOfDecks = nbrOfDecks
	return o
}

type OptFunc func(*Opts)

// Returns a default deck of Cards, see Func NewOpts for Defaults
// Default order:
// Ace of Diamonds, 2 of Diamonds ... King of Diamonds
// Ace of Hearts, 2 of Hearts ... King of Hearts
// Ace of Spade, 2 of Hearts ... King of Spade
// Ace of Clubs, 2 of Clubs ... King of Clubs
func New() []Card {
	return applyOpts(getDefaultCards(), NewOpts())
}

func NewWithOpts(opts Opts) []Card {
	return applyOpts(getDefaultCards(), opts)
}

func getDefaultCards() (cards []Card) {
	for _, rank := range ranks {
		for _, suit := range suits {
			cards = append(cards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}
	return
}

func getRankFilter(rank Rank) func(c Card) bool {
	return func(c Card) bool { return c.Rank == rank }
}

func getSuitFilter(suit Suit) func(c Card) bool {
	return func(c Card) bool { return c.Suit == suit }
}

func applyOpts(cards []Card, opts Opts) []Card {
	cards = applyJokers(opts.NbrOfJokers)(cards)
	cards = applyMultipleDecks(opts.NbrOfDecks)(cards)

	for _, rank := range opts.RanksToFilter {
		cards = applyFilter(getRankFilter(rank))(cards)
	}

	for _, suit := range opts.SuitsToFilter {
		cards = applyFilter(getSuitFilter(suit))(cards)
	}

	comp := getCompFunc(opts)
	sort.Slice(cards, comp(cards))

	if opts.WithShuffle {
		cards = applyShuffle(cards)
	}

	return cards
}

func getCompFunc(opts Opts) func(cards []Card) func(i, j int) bool {
	return func(cards []Card) func(i, j int) bool {
		return func(i, j int) bool {
			if opts.Reverse {
				return opts.rankFuncForSort(cards[i]) > opts.rankFuncForSort(cards[j])
			}
			return opts.rankFuncForSort(cards[i]) < opts.rankFuncForSort(cards[j])
		}
	}
}

func AbsRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

func applyJokers(nbr int) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < nbr; i++ {
			cards = append(cards, Card{Suit: Joker, Rank: Rank(i)})
		}
		return cards
	}
}

func applyMultipleDecks(nbrOfDecks int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var newDeck []Card
		for i := 0; i < nbrOfDecks; i++ {
			newDeck = append(newDeck, cards...)
		}
		return newDeck
	}
}

func applyFilter(filterFunc func(Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var newCards []Card
		for _, c := range cards {
			if !filterFunc(c) {
				newCards = append(newCards, c)
			}
		}
		return newCards
	}
}

func applyShuffle(cards []Card) []Card {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return cards
}
