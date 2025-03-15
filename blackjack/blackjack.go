package blackjack

import (
	"blackjack/deck"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
	score := 0
	for _, c := range h {
		if c.Rank == deck.Ace {
			tmp := score + 11
			if tmp > 21 {
				score += 1
			} else {
				score += 11
			}
		} else {
			score += min(int(c.Rank), 10)
		}
	}
	return score
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

func (p Player) DisplayHand(w io.Writer) int {
	p.IsDealer = false
	return p.DisplayVisibleHand(w)
}

// Writes to io.Writer. It returns the number of bytes written and any write error encountered
func (p Player) DisplayVisibleHand(w io.Writer) int {
	var hand = p.Hand
	if p.IsDealer {
		hand = p.Hand[:1]
	}

	var nbrBytes int
	for _, card := range hand {
		text := card.String() + "\n"
		n := DisplayMessage(w, text)
		nbrBytes += n
	}
	if p.IsDealer {
		text := "***HIDDEN***\n"
		n := DisplayMessage(w, text)
		nbrBytes += n

	}
	return nbrBytes
}

func (p Player) MakeAChoice(w io.Reader) (Choice, error) {
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

func ReadInput(r io.Reader) string {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Got unexpected Error %q", err)
	}
	return input[:len(input)-1]
}

func getChoice(w io.Reader) (Choice, error) {
	input := ReadInput(w)
	choice, err := ParseChoice(input)
	if err != nil && err != ErrInvalidInput {
		log.Fatal(err)
	}
	return choice, err
}

func ParseChoice(input string) (Choice, error) {
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
	DisplayMessage(w, msg)
}

func DisplayMessage(w io.Writer, msg string) int {
	n, err := w.Write([]byte(msg))
	if err != nil {
		log.Fatalf("Got unexpected Error %q", err)
	}
	return n
}

func DisplayGameAllVisible(w io.Writer, players []Player) error {
	players[len(players)-1].IsDealer = false
	return DisplayGame(w, players)
}

func DisplayGame(w io.Writer, players []Player) error {
	DisplayMessage(w, "Current state of the game:\n\n")
	for idx, player := range players[:len(players)-1] {
		fmt.Fprintf(w, "Player %d:\n", idx+1)
		player.DisplayVisibleHand(w)
	}

	dealer := players[len(players)-1]
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Dealer:\n")
	dealer.DisplayVisibleHand(w)
	fmt.Fprintf(w, "\n")

	return nil
}

func DisplayScore(w io.Writer, p Player) {
	score := p.Hand.Score()

	msg := fmt.Sprintf("Your current score is %d:\n", score)
	if p.IsDealer {
		msg = fmt.Sprintf("The dealers score is %d:\n", score)
	}

	DisplayMessage(w, msg)
}

// Hands every Player a Card, in order of the list of players
func DealOneRound(players []*Player, cards []deck.Card) ([]deck.Card, error) {
	var card deck.Card
	var err error
	for _, player := range players {
		card, cards, err = DealACard(cards)
		if err != nil {
			return cards, err
		}
		player.Hand = append(player.Hand, card)
	}
	return cards, nil
}

// Pops a card from cards, returns the card and the new set of cards
func DealACard(cards []deck.Card) (deck.Card, []deck.Card, error) {
	if len(cards) <= 0 {
		return deck.Card{}, cards, ErrEmptyDeck
	}
	card, cards := cards[len(cards)-1], cards[:len(cards)-1]
	return card, cards, nil
}

func DealCards(players []*Player) ([]deck.Card, error) {
	opts := deck.NewOpts().DoShuffle()
	cards := deck.NewWithOpts(opts)
	cards, err := DealOneRound(players, cards)
	if err != nil {
		log.Fatalf("error occured: %q", err)
	}
	cards, err = DealOneRound(players, cards)
	if err != nil {
		log.Fatalf("error occured: %q", err)
	}
	return cards, nil
}

func PlayRound(w io.Writer, r io.Reader, player *Player, cards []deck.Card) ([]deck.Card, error) {
	var card deck.Card
	for player.Hand.Score() <= 21 {
		DisplayMessage(w, "##############################\n")
		// DisplayGame(w, players)
		DisplayChoiceMessage(w)
		DisplayScore(w, *player)
		choice, err := player.MakeAChoice(r)
		if err == ErrInvalidInput {
			DisplayMessage(w, "Get the hell outta here\n")
			os.Exit(1)
			break
		}
		if choice == Stand {
			break
		}
		DisplayMessage(w, "Dealing...\n")
		time.Sleep(time.Second)

		card, cards, err = DealACard(cards)
		if err != nil {
			if err == ErrEmptyDeck {
				log.Println("Sorry Deck is Empty")
				break
			} else {
				log.Fatal(err)
			}
		}
		player.Hand = append(player.Hand, card)
	}
	if player.Hand.Score() > 21 {
		DisplayMessage(w, "You are over 21\n")
	}
	return cards, nil

}
