package blackjack

import (
	"blackjackAI/deck"
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

const (
	ErrEmptyDeck    = BlackjackError("Empty Deck")
	ErrInvalidInput = BlackjackError("Invalid Input")
)

type BlackjackError string

func (e BlackjackError) Error() string {
	return string(e)
}

type Choice uint8

const (
	_ Choice = iota
	Hit
	Stand
)

type Hand []deck.Card

func (h Hand) Score() int {
	score := h.minScore()
	for _, c := range h {
		if c.Rank == deck.Ace {
			if score+10 <= 21 {
				score += 10
			}
		}
	}
	return score
}

func (h Hand) minScore() (score int) {
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return
}

func (h Hand) String() string {
	var strs []string
	for _, c := range h {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Player struct {
	Hand     Hand
	IsDealer bool
}

func (p Player) displayHand(w io.Writer) {
	p.IsDealer = false
	p.DisplayVisibleHand(w)
	return
}

// Writes to io.Writer. It returns the number of bytes written and any write error encountered
func (p Player) DisplayVisibleHand(w io.Writer) {
	if p.IsDealer {
		fmt.Fprintf(w, p.Hand[:1].String())
	}
	fmt.Fprintf(w, p.Hand[:1].String())
	return
}

func (p Player) MakeChoice(w io.Reader) (Choice, error) {
	var choice Choice
	var err error
	count := 0
	for count < 5 {
		if p.IsDealer {
			choice = p.getDealerChoice()
		} else {
			choice, err = getChoice(w)
		}
		if err == nil {
			break
		}
		count++
	}
	if err == ErrInvalidInput {
		return 0, err
	}
	return choice, nil
}

func (p Player) getDealerChoice() Choice {
	score := p.Hand.Score()
	withAce := func(h Hand) (hasAce bool) {
		for _, c := range h {
			if c.Rank == deck.Ace {
				return true
			}
		}
		return
	}(p.Hand)

	if score < 17 || (score == 17 && withAce) {
		return Hit
	}
	return Stand
}

func readInput(r io.Reader) string {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Got unexpected Error %q", err)
	}
	return input[:len(input)-1]
}

func getChoice(w io.Reader) (Choice, error) {
	input := readInput(w)
	choice, err := parseChoice(input)
	if err != nil && err != ErrInvalidInput {
		log.Fatal(err)
	}
	return choice, err
}

func parseChoice(input string) (Choice, error) {
	if input == "1" {
		return Hit, nil
	} else if input == "2" {
		return Stand, nil
	}
	return 0, ErrInvalidInput
}

func DisplayChoiceMessage(w io.Writer) {
	msg := `
Please make your choice:
1 - hit
2 - stand
`
	fmt.Printf("%s\n", msg)
}

type Options struct {
	Hands int
	Decks int
}

type Game struct {
	opts   Options
	deck   []deck.Card
	State  State
	player AI
	dealer Player
}

type AI interface {
	MakeChoice() Choice
}

func (g Game) Play(ai AI) {

	ai.MakeChoice()
	return
}

func New(opts Options) Game {
	return Game{}
}
