package game

import (
	"fmt"
	"math/rand"
	"sort"
)

const CardTypes = 13
const InitialCards = 10

type Card int

type Cards []Card

func (cs Cards) Len() int           { return len(cs) }
func (cs Cards) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs Cards) Less(i, j int) bool { return cs[i] < cs[j] }

type Player struct {
	Name  string
	Cards Cards
	Phase int
	Out   bool
}

type Game struct {
	Players  []Player
	OutCards []Cards
	Turn     int
	Trash    Card
}

type PlayerType int

const (
	Human PlayerType = 0
	AI
)

func SetUp(playerCount int, turn int) *Game {
	var players []Player
	for i := 0; i < playerCount; i++ {
		players = append(players, Player{
			fmt.Sprintf("Player %v", i+1),
			serveN(InitialCards),
			1, false,
		})
	}

	return &Game{
		players,
		[]Cards{},
		turn, serve(),
	}
}

func (g *Game) Draw(fromTrash bool) {
	var card Card
	if fromTrash {
		card = g.Trash
		g.Trash = -1
	} else {
		card = serve()
	}

	g.Players[g.Turn].Cards = append(g.Players[g.Turn].Cards, card)
	sort.Sort(g.Players[g.Turn].Cards)
}

func serveN(n int) (cards Cards) {
	for i := 0; i < n; i++ {
		cards = append(cards, serve())
	}
	sort.Sort(cards)
	return
}

func serve() Card {
	return Card((rand.Int() % CardTypes) + 1)
}

func (g *Game) Drop(idx int) {
	g.Trash = g.Players[g.Turn].Cards[idx]

	g.Players[g.Turn].Cards = append(
		g.Players[g.Turn].Cards[:idx],
		g.Players[g.Turn].Cards[idx+1:]...)

	g.Turn = (g.Turn + 1) % len(g.Players)
}

func (g *Game) ComeOut(idxSeq [][]int) {
	if len(idxSeq) == 0 {
		return // TODO error
	}
	// TODO check duplicates

	for _, idxs := range idxSeq {
		cards := Cards{}
		for _, idx := range idxs {
			cards = append(cards, g.Players[g.Turn].Cards[idx])
		}
		g.OutCards = append(g.OutCards, cards)
	}

	var idxs []int
	for _, xs := range idxSeq {
		idxs = append(idxs, xs...)
	}
	sort.Ints(idxs)
	newCards := Cards{}
	j := 0
	for _, i := range idxs {
		newCards = append(newCards, g.Players[g.Turn].Cards[j:i]...)
		j = i + 1
	}
	newCards = append(newCards, g.Players[g.Turn].Cards[j:]...)

	g.Players[g.Turn].Cards = newCards
	g.Players[g.Turn].Out = true
}

func (g *Game) Append(card, sequence int, left bool) {
	c := g.Players[g.Turn].Cards[card]
	g.Players[g.Turn].Cards = append(
		g.Players[g.Turn].Cards[:card],
		g.Players[g.Turn].Cards[card+1:]...)

	if left {
		g.OutCards[sequence] = append(Cards{c}, g.OutCards[sequence]...)
	} else {
		g.OutCards[sequence] = append(g.OutCards[sequence], c)
	}
}

func (g *Game) IsDone() bool {
	for _, p := range g.Players {
		if len(p.Cards) == 0 {
			return true
		}
	}

	return false
}

func (g *Game) GetWinner() int {
	for i, p := range g.Players {
		if len(p.Cards) == 0 {
			return i
		}
	}

	return -1
}
