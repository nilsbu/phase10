package game

import (
	"errors"
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

func (g *Game) Drop(idx int) error {
	if idx < 0 || idx >= len(g.Players[g.Turn].Cards) {
		return fmt.Errorf("index out of bounds: %v", idx)
	}

	g.Trash = g.Players[g.Turn].Cards[idx]

	g.Players[g.Turn].Cards = append(
		g.Players[g.Turn].Cards[:idx],
		g.Players[g.Turn].Cards[idx+1:]...)

	g.Turn = (g.Turn + 1) % len(g.Players)

	return nil
}

func (g *Game) ComeOut(idxSeq [][]int) error {
	if len(idxSeq) == 0 {
		return errors.New("no cards")
	}
	if containsDuplicates(idxSeq) {
		return errors.New("cards contain duplicates")
	}
	if g.Players[g.Turn].Out {
		return errors.New("player is already out")
	}

	var cardss []Cards
	var seqs []Sequence

	for _, idxs := range idxSeq {
		cards := Cards{}
		for _, idx := range idxs {
			cards = append(cards, g.Players[g.Turn].Cards[idx])
		}

		seq := validate(cards)
		if seq.Type == Invalid || seq.Type == Ambiguous {
			return errors.New("invalid cards")
		}

		cardss = append(cardss, cards)
		seqs = append(seqs, seq)
	}

	if !isPhaseFulfilled(seqs, g.Players[g.Turn].Phase) {
		return errors.New("phase is not fulfilled")
	}

	for _, cards := range cardss {
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

	return nil
}

func containsDuplicates(idxSeq [][]int) bool {
	var idxs []int
	for _, seq := range idxSeq {
		idxs = append(idxs, seq...)
	}
	sort.Ints(idxs)
	for i := 0; i < len(idxs)-1; i++ {
		if idxs[i] == idxs[i+1] {
			return true
		}
	}
	return false
}

func (g *Game) Append(card, sequence int, left bool) error {
	if card < 0 || card >= len(g.Players[g.Turn].Cards) {
		return errors.New("card index out of bounds")
	}
	if sequence < 0 || sequence >= len(g.OutCards) {
		return errors.New("sequence index out of bounds")
	}
	if !g.Players[g.Turn].Out {
		return errors.New("player is not out")
	}

	c := g.Players[g.Turn].Cards[card]

	var newSeq Cards
	if left {
		newSeq = append(Cards{c}, g.OutCards[sequence]...)
	} else {
		newSeq = append(g.OutCards[sequence], c)
	}

	if seq := validate(newSeq); seq.Type == Invalid || seq.Type == Ambiguous {
		return errors.New("sequence invalid")
	}

	g.OutCards[sequence] = newSeq

	g.Players[g.Turn].Cards = append(
		g.Players[g.Turn].Cards[:card],
		g.Players[g.Turn].Cards[card+1:]...)

	return nil
}

func (g *Game) IsDone() bool {
	allOut := true
	for _, p := range g.Players {
		if len(p.Cards) == 0 {
			return true
		}
		if !p.Out {
			allOut = false
		}
	}

	return allOut
}
