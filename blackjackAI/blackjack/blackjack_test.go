package blackjack

import (
	"blackjackAI/deck"
	"bytes"
	"reflect"
	"testing"
)

func TestScore(t *testing.T) {
	hand := Hand(
		[]deck.Card{
			{Suit: deck.Club, Rank: deck.Ace},
			{Suit: deck.Club, Rank: deck.Ten},
			{Suit: deck.Club, Rank: deck.Ace},
		})
	want := 12
	got := hand.Score()
	if got != want {
		t.Errorf("hand.Score() = \"%v\"; want \"%v\"", got, want)
	}
}

func TestPlayersTurn(t *testing.T) {
	t.Run("Test DisplayMessage", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		msg := "Test message"
		DisplayMessage(buffer, msg)

		got := buffer.String()
		want := msg
		if got != want {
			t.Errorf("DisplayMessage() = \"%v\"; want \"%v\"", got, want)
		}
	})
	t.Run("Accept 1 or 2", func(t *testing.T) {
		testCases := []struct {
			input string
			want  Choice
		}{
			{input: "1", want: Hit},
			{input: "2", want: Stand},
		}
		for _, tc := range testCases {
			got, err := ParseChoice(tc.input)

			assertNoError(t, err)
			if got != tc.want {
				t.Errorf("() = \"%v\"; want \"%v\"", got, tc.want)
			}
		}
	})
	t.Run("Throw Error if not 1 or 2", func(t *testing.T) {
		input := "totaler Quatsch"
		_, err := ParseChoice(input)
		assertError(t, err, ErrInvalidInput)
	})
}

func TestDisplay(t *testing.T) {
	t.Run("Test Display of Game", func(t *testing.T) {
		buffer := bytes.Buffer{}
		players := []Player{
			{Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
				deck.Card{Suit: deck.Heart, Rank: deck.Nine},
			}},
			{Hand: Hand{
				deck.Card{Suit: deck.Diamond, Rank: deck.Ten},
				deck.Card{Suit: deck.Spade, Rank: deck.Two},
			}},
		}
		DisplayGame(&buffer, players)

		got := buffer.String()
		want := `Current state of the game:

Player 1:
King of Clubs
Nine of Hearts

Dealer:
Ten of Diamonds
Two of Spades

`
		if got != want {
			t.Errorf("buffer.String() = \"%v\"; want \"%v\"", got, want)
		}

	})
	t.Run("Test DisplayVisibleHand", func(t *testing.T) {
		buffer := bytes.Buffer{}
		player := Player{
			Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
				deck.Card{Suit: deck.Heart, Rank: deck.Nine},
			},
		}
		player.DisplayVisibleHand(&buffer)

		got := buffer.String()
		want := `King of Clubs
Nine of Hearts
`
		if got != want {
			t.Errorf("buffer.String() = \"%v\"; want \"%v\"", got, want)
		}
	})

	t.Run("Test DisplayVisibleHand Dealer One Card", func(t *testing.T) {
		buffer := bytes.Buffer{}
		player := Player{
			Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
				deck.Card{Suit: deck.Heart, Rank: deck.Nine},
			},
			IsDealer: true,
		}
		player.DisplayVisibleHand(&buffer)

		got := buffer.String()
		want := `King of Clubs
***HIDDEN***
`
		if got != want {
			t.Errorf("buffer.String() = \"%v\"; want \"%v\"", got, want)
		}
	})
	t.Run("Test DisplayHand Dealer", func(t *testing.T) {
		buffer := bytes.Buffer{}
		player := Player{
			Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
				deck.Card{Suit: deck.Heart, Rank: deck.Nine},
			},
			IsDealer: true,
		}
		player.DisplayHand(&buffer)

		gotP := player.IsDealer
		got := buffer.String()
		want := `King of Clubs
Nine of Hearts
`
		if !gotP {
			t.Errorf("Wanted him to stay the dealer")
		}
		if got != want {
			t.Errorf("buffer.String() = \"%v\"; want \"%v\"", got, want)
		}
	})
}

func TestDealing(t *testing.T) {
	t.Run("Test DealACard reduces deck size by one", func(t *testing.T) {
		cards := deck.New()
		oldLength := len(cards)
		_, cards, err := DealACard(cards)
		newLength := len(cards)
		want := oldLength - 1

		assertNoError(t, err)

		if newLength != want {
			t.Errorf("Want newLength to be %d; got %d", want, newLength)
		}
	})

	t.Run("Test DealACard takes from top of Deck", func(t *testing.T) {
		cards := deck.New() // top of Deck is King of Clubs
		got, _, err := DealACard(cards)
		want := deck.Card{Suit: deck.Club, Rank: deck.King}

		assertNoError(t, err)
		if got != want {
			t.Errorf("Got %v; want %v", got, want)
		}
	})

	t.Run("ErrEmptyDeck ", func(t *testing.T) {
		cards := []deck.Card{}
		_, _, err := DealACard(cards)
		assertError(t, err, ErrEmptyDeck)
	})

	t.Run("Test Deal", func(t *testing.T) {
		cards := deck.New()
		players := []*Player{{}, {}}
		DealOneRound(players, cards)

		got := []Player{*players[0], *players[1]}
		want := []Player{
			{Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
			}},
			{Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.Queen},
			}},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v; want %v", got, want)
		}
	})

	t.Run("Test Deal return EmptyDeck Error", func(t *testing.T) {
		cards := []deck.Card{}
		_, err := DealOneRound([]*Player{{}, {}}, cards)
		assertError(t, err, ErrEmptyDeck)
	})

	t.Run("Test Deal stops Dealing when ErrEmptyDeck encountered", func(t *testing.T) {
		cards := []deck.Card{
			{Suit: deck.Club, Rank: deck.King},
		}
		players := []*Player{{}, {}}
		cards, err := DealOneRound(players, cards)

		assertError(t, err, ErrEmptyDeck)

		got := []Player{*players[0], *players[1]}
		want := []Player{
			{Hand: Hand{
				deck.Card{Suit: deck.Club, Rank: deck.King},
			}},
			{},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v; want %v", got, want)
		}

	})
}

func assertError(t testing.TB, gotErr, wantErr error) {
	t.Helper()
	if gotErr != wantErr {
		t.Errorf("Got error %q; want error %q", gotErr, wantErr)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Got an error: %q; didn't want one", err)
	}
}
