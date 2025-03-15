package main

import (
	bj "blackjack"
	"os"
	"time"
)

func main() {
	players := []*bj.Player{
		{},
		{IsDealer: true},
	}

	bj.DisplayMessage(os.Stdout, "Dealing...\n")
	cards, err := bj.DealCards(players)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	cards, err = bj.PlayRound(os.Stdout, os.Stdin, players[0], cards)
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	bj.DisplayMessage(os.Stdout, "##############################\n")
	bj.DisplayMessage(os.Stdout, "Dealer making choices\n")
	dealer := players[len(players)-1]
	cards, err = bj.PlayRound(os.Stdout, os.Stdin, dealer, cards)
	if err != nil {
		panic(err)
	}

	// bj.DisplayMessage(os.Stdout, "##############################\n")
	// bj.DisplayMessage(os.Stdout, "End of game\n")
	// bj.DisplayGameAllVisible(os.Stdout, []bj.Player{player, dealer})
	// bj.DisplayScore(os.Stdout, player)
	// bj.DisplayScore(os.Stdout, dealer)

	// if player.Hand.Score() > 21 {
	// 	bj.DisplayMessage(os.Stdout, "You lose\n")
	// 	return
	// }

	// if player.Hand.Score() > dealer.Hand.Score() {
	// 	bj.DisplayMessage(os.Stdout, "You won\n")
	// 	return
	// }
	// bj.DisplayMessage(os.Stdout, "You lose\n")
	return
}
