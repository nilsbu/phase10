package actor

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/phase10/pkg/game"
)

func TestSplitAddString(t *testing.T) {
	cs := []struct {
		str  string
		cmds []addCmd
	}{
		{
			"0<2;2>22;9>7;",
			[]addCmd{
				{0, 2, true}, {2, 22, false}, {9, 7, false},
			},
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			cmds := splitAddString(c.str)
			if !reflect.DeepEqual(c.cmds, cmds) {
				t.Errorf("expected cmds:\n%v\nbut got:\n%v", c.cmds, cmds)
			}
		})
	}
}

func TestSplitCOString(t *testing.T) {
	cs := []struct {
		name string
		str  string
		cmds [][]int
	}{
		{
			"not done",
			"1,2,3;4,66,9;",
			[][]int{
				{1, 2, 3}, {4, 66, 9},
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			cmds := splitCOString(c.str)
			if !reflect.DeepEqual(c.cmds, cmds) {
				t.Errorf("expected cmds:\n%v\nbut got:\n%v", c.cmds, cmds)
			}
		})
	}
}

func TestFindCards(t *testing.T) {
	cs := []struct {
		want []int
		has  game.Cards
		idxs []int
		err  error
	}{
		{
			[]int{1, 4, 4, 8},
			game.Cards{1, 1, 2, 4, 4, 4, 7, 8, 9},
			[]int{0, 3, 4, 7},
			nil,
		},
		{
			[]int{1, 4, 4, 8, 8},
			game.Cards{1, 1, 2, 4, 4, 4, 7, 8, 9},
			[]int{0, 3, 4, 7},
			errors.New("no 2nd 8"),
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			idxs, err := findCards(c.want, c.has)
			if (err == nil) != (c.err == nil) {
				t.Errorf("error unexpected: got '%v', expected '%v'", err, c.err)
			}
			if err == nil && !reflect.DeepEqual(c.idxs, idxs) {
				t.Errorf("expected cmds:\n%v\nbut got:\n%v", c.idxs, idxs)
			}
		})
	}
}