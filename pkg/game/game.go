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

type CardSequence []Cards

func (cs CardSequence) Len() int      { return len(cs) }
func (cs CardSequence) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs CardSequence) Less(i, j int) bool {
	vi, vj := Validate(cs[i]), Validate(cs[j])

	if vi.Type > vj.Type {
		return true
	} else if vi.Type < vj.Type {
		return false
	} else if vi.Type == Invalid {
		return vi.N > vj.N
	} else if vi.Type == Kind {
		var ni, nj Card
		for _, c := range cs[i] {
			ni = c
			if c != 13 {
				break
			}
		}
		for _, c := range cs[j] {
			nj = c
			if c != 13 {
				break
			}
		}
		if ni != nj {
			return ni < nj
		}
		return vi.N > vj.N
	} else {
		fi, fj := firstCardOfStraight(cs[i]), firstCardOfStraight(cs[j])
		if fi > fj {
			return false
		} else if fi < fj {
			return true
		}
		return vi.N > vj.N
	}
}

func firstCardOfStraight(cs Cards) Card {
	f := Card(-1)
	for i, c := range cs {
		if c != 13 {
			f = c - Card(i)
			break
		}
	}
	return f
}

type Player struct {
	Name  string
	Cards Cards
	Phase int
	Out   bool
}

type Game struct {
	Players  []Player
	OutCards CardSequence
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

func (g *Game) Drop(card Card) error {
	if err := g.removeCards([]Cards{{card}}); err != nil {
		return err
	}

	g.Trash = card

	g.Turn = (g.Turn + 1) % len(g.Players)

	sort.Sort(g.OutCards)

	return nil
}

func (g *Game) removeCards(toRemove []Cards) error {
	remaining := Cards{}
	for _, card := range g.Players[g.Turn].Cards {
		remaining = append(remaining, card)
	}

	for _, cards := range toRemove {
		for _, card := range cards {
			idx, err := findCard(remaining, card)
			if err != nil {
				return err
			}
			remaining = append(remaining[:idx], remaining[idx+1:]...)
		}
	}

	g.Players[g.Turn].Cards = remaining
	return nil
}

func findCard(cards Cards, card Card) (idx int, err error) {
	var c Card
	for idx, c = range cards {
		if card == c {
			return
		}
	}

	return -1, fmt.Errorf("card %v not found", card)
}

func (g *Game) ComeOut(cardss []Cards) error {
	if g.Players[g.Turn].Out {
		return errors.New("player is already out")
	}

	var seqs []Sequence
	for _, cards := range cardss {
		seqs = append(seqs, Validate(cards))
	}

	if !isPhaseFulfilled(seqs, g.Players[g.Turn].Phase) {
		return errors.New("phase is not fulfilled")
	}

	if err := g.removeCards(cardss); err != nil {
		return err
	}

	for _, cards := range cardss {
		g.OutCards = append(g.OutCards, cards)
	}

	g.Players[g.Turn].Out = true

	return nil
}

func (g *Game) Append(card Card, sequence int, left bool) error {
	if sequence < 0 || sequence >= len(g.OutCards) {
		return errors.New("sequence index out of bounds")
	}
	if !g.Players[g.Turn].Out {
		return errors.New("player is not out")
	}

	var newSeq Cards
	if left {
		newSeq = append(Cards{card}, g.OutCards[sequence]...)
	} else {
		newSeq = append(g.OutCards[sequence], card)
	}

	if seq := Validate(newSeq); seq.Type == Invalid || seq.Type == Ambiguous {
		return errors.New("sequence invalid")
	}

	if err := g.removeCards([]Cards{{card}}); err != nil {
		return err
	}

	g.OutCards[sequence] = newSeq

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

func (g *Game) GetWinner() int {
	if !g.IsDone() {
		return -1
	}
	winner := -1
	for i, p := range g.Players {
		if p.Phase == 10 && p.Out {
			if winner == -1 {
				winner = i
			} else if len(p.Cards) == 0 && len(g.Players[winner].Cards) > 0 {
				winner = i
			} else {
				return -1
			}
		}
	}

	return winner
}

func (g *Game) NextRound() {
	for i := range g.Players {
		g.Players[i].Cards = serveN(InitialCards)

		if g.Players[i].Out {
			g.Players[i].Out = false
			g.Players[i].Phase++
		}
	}

	g.OutCards = []Cards{}
	g.Trash = serve()
}
