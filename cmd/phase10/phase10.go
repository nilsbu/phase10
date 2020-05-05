package main

import (
	"fmt"
	"math/rand"
	"time"

	a "github.com/nilsbu/phase10/pkg/actor"
	g "github.com/nilsbu/phase10/pkg/game"
)

func main() {
	rand.Seed(time.Now().Unix())
	game := g.SetUp(2, 0)

	for {
		actor := &a.Human{}
		if err := actor.Play(game); err != nil {
			fmt.Println(err)
			return
		}
		if game.GetWinner() > -1 {
			fmt.Printf("the winner is %v\n",
				game.Players[game.GetWinner()].Name)
			return
		}
		if game.IsDone() {
			fmt.Println("round finished")
			game.NextRound()
		}
	}
}
