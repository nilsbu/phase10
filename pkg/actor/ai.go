package actor

import (
	"fmt"
	"time"

	"github.com/nilsbu/phase10/pkg/display"
	"github.com/nilsbu/phase10/pkg/game"
)

type AI struct{}

func (h *AI) Play(g *game.Game) error {
	if err := h.draw(g); err != nil {
		return err
	}
	fmt.Println(">> drawn")
	fmt.Println(display.PrintGame(g, g.Turn))
	time.Sleep(1 * time.Second)

	if err := h.put(g); err != nil {
		return err
	}
	fmt.Println(">> put")
	if len(g.Players[g.Turn].Cards) == 0 {
		g.Turn = (g.Turn + 1) % len(g.Players)
		return nil
	}
	fmt.Println(display.PrintGame(g, g.Turn))
	time.Sleep(1 * time.Second)

	if err := h.drop(g); err != nil {
		return err
	}
	fmt.Println(">> drop")
	return nil
}

func (h *AI) draw(g *game.Game) error {
	g.Draw(g.Trash == 13)
	return nil
}

func (h *AI) put(g *game.Game) error {
	if err := comeOut(g); err != nil {
		return err
	}

	if err := appendCards(g); err != nil {
		return err
	}

	return nil
}

func comeOut(g *game.Game) error {
	if g.Players[g.Turn].Out {
		return nil
	}

	seqs, _ := game.GetPhaseSequences(g.Players[g.Turn].Phase)
	if cardss, ok := comeOutR(g.Players[g.Turn].Cards, seqs); ok {

		if err := g.ComeOut(cardss); err != nil {
			return err
		}
	}

	return nil
}

func comeOutR(cards game.Cards, seqs []game.Sequence) ([]game.Cards, bool) {
	if len(seqs) == 0 {
		return []game.Cards{}, true
	}
	switch seqs[0].Type {
	case game.Kind:
		for i, card := range cards {
			idxs := game.Cards{} // TODO change name
			for j := i; j < len(cards); j++ {
				if cards[j] != card && cards[j] != 13 {
					continue
				}
				idxs = append(idxs, cards[j])
				if len(idxs) >= seqs[0].N {
					if oidxs, ok := comeOutR(removeCards(cards, idxs), seqs[1:]); ok {
						return append([]game.Cards{idxs}, oidxs...), true
					}
					break
				}
			}
		}

	case game.Straight:
		nums, jokers := splitOffJokers(cards)

		for start := 1; start <= 13-seqs[0].N; start++ {
			idxs := game.Cards{} // TODO change name
			n, j := 0, 0
			for i := start; i <= 12; i++ {
				for {
					if n < len(nums) && nums[n] == game.Card(start)+game.Card(len(idxs)) {
						idxs = append(idxs, nums[n])
						n++
						break
					} else if j < len(jokers) {
						idxs = append(idxs, 13)
						j++
						break
					} else {
						if n >= len(nums) && j >= len(jokers) {
							break
						}
						n++
					}
				}

				if len(idxs) >= seqs[0].N {
					if oidxs, ok := comeOutR(removeCards(cards, idxs), seqs[1:]); ok {
						return append([]game.Cards{idxs}, oidxs...), true
					}
					break
				}
			}
		}
	}

	return []game.Cards{}, false
}

func removeCards(cards game.Cards, toRemove game.Cards) game.Cards {
	out := game.Cards{}
	i := 0
	for _, card := range cards {
		if i < len(toRemove) && card == toRemove[i] {
			i++
		} else {
			out = append(out, card)
		}
	}
	return out
}

func splitOffJokers(cards game.Cards) (game.Cards, game.Cards) {
	for i, card := range cards {
		if card == 13 {
			return cards[:i], cards[i:]
		}
	}
	return cards, game.Cards{}
}

func appendCards(g *game.Game) error {
	if !g.Players[g.Turn].Out {
		return nil
	}

	for i := range g.OutCards {
		for j := 0; j < len(g.Players[g.Turn].Cards); {
			fmt.Println(i, j)
			if err := g.Append(g.Players[g.Turn].Cards[j], i, false); err == nil {
				j = 0
				continue
			}
			if err := g.Append(g.Players[g.Turn].Cards[j], i, true); err == nil {
				j = 0
				continue
			}
			j++
		}
	}
	return nil
}

func (h *AI) drop(g *game.Game) error {
	return g.Drop(g.Players[g.Turn].Cards[0])
}
