package actor

import (
	"reflect"
	"testing"

	"github.com/nilsbu/phase10/pkg/game"
)

func TestComeOutR(t *testing.T) {
	cs := []struct {
		cards game.Cards
		seqs  []game.Sequence
		out   []game.Cards
		ok    bool
	}{
		{
			game.Cards{},
			[]game.Sequence{},
			[]game.Cards{}, true,
		},
		{
			game.Cards{2, 3, 13},
			[]game.Sequence{},
			[]game.Cards{}, true,
		},
		{
			game.Cards{1, 3, 4, 4, 4, 5},
			[]game.Sequence{{Type: game.Kind, N: 3}},
			[]game.Cards{{4, 4, 4}}, true,
		},
		{
			game.Cards{1, 3, 4, 4, 4, 5},
			[]game.Sequence{{Type: game.Kind, N: 4}},
			[]game.Cards{}, false,
		},
		{
			game.Cards{1, 3, 4, 5, 6, 13, 13},
			[]game.Sequence{{Type: game.Kind, N: 3}},
			[]game.Cards{{1, 13, 13}}, true,
		},
		{
			game.Cards{1, 4, 5, 5, 6, 8},
			[]game.Sequence{{Type: game.Straight, N: 3}},
			[]game.Cards{{4, 5, 6}}, true,
		},
		{
			game.Cards{2, 4, 13},
			[]game.Sequence{{Type: game.Straight, N: 3}},
			[]game.Cards{{2, 13, 4}}, true,
		},
		{
			game.Cards{6, 7, 10},
			[]game.Sequence{{Type: game.Straight, N: 4}},
			[]game.Cards{}, false,
		},
		{
			game.Cards{10, 11, 12, 13},
			[]game.Sequence{{Type: game.Straight, N: 4}},
			[]game.Cards{{13, 10, 11, 12}}, true,
		},
		{
			game.Cards{12},
			[]game.Sequence{{Type: game.Straight, N: 1}},
			[]game.Cards{{12}}, true,
		},
		{
			game.Cards{2, 3, 3, 3, 3, 4},
			[]game.Sequence{{Type: game.Straight, N: 3}, {Type: game.Kind, N: 3}},
			[]game.Cards{{2, 3, 4}, {3, 3, 3}}, true,
		},
		{
			game.Cards{2, 3, 3, 3, 4},
			[]game.Sequence{{Type: game.Straight, N: 3}, {Type: game.Kind, N: 3}},
			[]game.Cards{}, false,
		},
		{
			game.Cards{2, 2, 2, 3, 4, 5},
			[]game.Sequence{{Type: game.Straight, N: 3}, {Type: game.Kind, N: 3}},
			[]game.Cards{{3, 4, 5}, {2, 2, 2}}, true,
		},
		{
			game.Cards{2, 2, 2, 3, 4, 6, 6, 6},
			[]game.Sequence{{Type: game.Kind, N: 3}, {Type: game.Straight, N: 3}},
			[]game.Cards{{6, 6, 6}, {2, 3, 4}}, true,
		},
		{
			game.Cards{2, 2, 2, 3, 4, 6, 6},
			[]game.Sequence{{Type: game.Kind, N: 3}, {Type: game.Straight, N: 3}},
			[]game.Cards{}, false,
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			out, ok := comeOutR(c.cards, c.seqs)
			if c.ok != ok {
				t.Errorf("expected '%v' but got '%v'", c.ok, ok)
			}
			if !reflect.DeepEqual(c.out, out) {
				t.Errorf("expected '%v' but got %v", c.out, out)
			}
		})
	}
}
