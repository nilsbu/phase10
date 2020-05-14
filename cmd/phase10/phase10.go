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
		var actor a.Actor
		if game.Turn == len(game.Players)-1 {
			actor = &a.Human{}
		} else {
			actor = &a.AI{}
		}
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
			out := ""
			for _, player := range game.Players {
				if player.Out {
					out += player.Name + ", "
				}
			}
			fmt.Printf("Players that came out: %v\n", out[:len(out)-2])
			fmt.Println("====== round finished ======")

			game.NextRound()
		}
	}
}
