package actor

import (
	"reflect"
	"testing"
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
