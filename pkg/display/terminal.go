package display

import (
	"fmt"
	"strconv"

	g "github.com/nilsbu/phase10/pkg/game"
)

const separator = "--------"

func PrintGame(game *g.Game, visible int) string {
	out := ""
	for i, p := range game.Players {
		out += printPlayer(p, i == visible)
		if i < len(game.Players)-1 {
			out += "\n"
		}
	}
	out += fmt.Sprintf("%v\n%v%v\n{%v,%v} - %v's turn",
		separator, printOutCards(game.OutCards), separator,
		PrintCard(-1, false), // any hidden card
		PrintCard(game.Trash, true),
		game.Players[game.Turn].Name)

	return out
}

func printPlayer(p g.Player, visible bool) string {
	out := ""
	if p.Out {
		out = "++"
	}
	return fmt.Sprintf(
		"%v (%v%v, %v)\n%v\n",
		p.Name, p.Phase, out, getPhaseRequirements(p.Phase),
		printCards(p.Cards, visible),
	)
}

func getPhaseRequirements(phase int) string {
	seqs, _ := g.GetPhaseSequences(phase)

	s := ""
	for _, seq := range seqs {
		var k string
		switch seq.Type {
		case g.Kind:
			k = "K"
		case g.Straight:
			k = "S"
		default:
			k = "?"
		}
		s += fmt.Sprintf("%v%v", k, seq.N)
	}
	return s
}

func printOutCards(outCards []g.Cards) string {
	s := ""
	for i, cards := range outCards {
		s += fmt.Sprintf(
			"%v: %v\n",
			i+1, printCards(cards, true),
		)
	}
	return s
}

func printCards(cards g.Cards, visible bool) string {
	s := "{"
	for i, c := range cards {
		s += PrintCard(c, visible)
		if i < len(cards)-1 {
			s += ","
		}
	}

	return s + "}"
}

func PrintCard(c g.Card, visible bool) string {
	switch {
	case !visible:
		return "X"
	case c >= 0 && c < 10:
		return strconv.Itoa(int(c))
	case c >= 10 && c <= 12:
		return string(c - 10 + 'a')
	case c == 13:
		return "J"
	default:
		return "?"
	}
}
