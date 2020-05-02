package display

import (
	"fmt"
	"testing"

	g "github.com/nilsbu/phase10/pkg/game"
)

func TestPrintPlayer(t *testing.T) {
	cs := []struct {
		name     string
		player   g.Player
		visible  bool
		expected string
	}{
		{
			"hidden",
			g.Player{
				Name:  "Player 1",
				Cards: g.Cards{1, 2, 3},
				Phase: 1,
				Out:   false,
			},
			false,
			"Player 1 (1)\n{X,X,X}\n",
		},
		{
			"shown",
			g.Player{
				Name:  "Player 1",
				Cards: g.Cards{1, 2, 11, 13},
				Phase: 1,
				Out:   false,
			},
			true,
			"Player 1 (1)\n{1,2,b,J}\n",
		},
		{
			"hidden out",
			g.Player{
				Name:  "Player 1",
				Cards: g.Cards{1, 2, 11, 13},
				Phase: 1,
				Out:   true,
			},
			false,
			"Player 1 (1++)\n{X,X,X,X}\n",
		},
		{
			"broken card",
			g.Player{
				Name:  "P1",
				Cards: g.Cards{-1},
				Phase: 1,
				Out:   false,
			},
			true,
			"P1 (1)\n{?}\n",
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			s := printPlayer(c.player, c.visible)
			if s != c.expected {
				t.Errorf("expected:\n%v\nbut got:\n%v", c.expected, s)
			}
		})
	}
}

func TestPrintOutCards(t *testing.T) {
	cs := []struct {
		name     string
		outCards []g.Cards
		expected string
	}{
		{
			"empty",
			[]g.Cards{},
			"",
		},
		{
			"1 sequence",
			[]g.Cards{{9, 10, 11, 12}},
			"1: {9,a,b,c}\n",
		},
		{
			"2 sequence2",
			[]g.Cards{{9, 10, 11, 12}, {5, 5, 13, 5, 5}},
			"1: {9,a,b,c}\n2: {5,5,J,5,5}\n",
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			s := printOutCards(c.outCards)
			if s != c.expected {
				t.Errorf("expected:\n%v\nbut got:\n%v", c.expected, s)
			}
		})
	}
}

func TestPrintGame(t *testing.T) {
	p1 := g.Player{Name: "P1",
		Cards: g.Cards{1, 1, 1},
		Phase: 4, Out: true}

	p2 := g.Player{Name: "P2",
		Cards: g.Cards{10, 11, 13},
		Phase: 10, Out: false}

	out := []g.Cards{
		{1, 1, 1}, {10, 10, 10}, {10, 13, 12},
	}

	cs := []struct {
		name     string
		game     *g.Game
		visible  int
		expected string
	}{
		{
			"game 1",
			&g.Game{
				Players:  []g.Player{p1, p2},
				OutCards: out,
				Turn:     1, Trash: 11,
			},
			0,
			fmt.Sprintf(
				"%v\n%v--------\n%v--------\n{X,b} - P2's turn",
				printPlayer(p1, true), printPlayer(p2, false),
				printOutCards(out)),
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			s := PrintGame(c.game, c.visible)
			if s != c.expected {
				t.Errorf("expected:\n%v\nbut got:\n%v", c.expected, s)
			}
		})
	}
}
