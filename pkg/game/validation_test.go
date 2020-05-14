package game

import "testing"

func TestValidate(t *testing.T) {
	cs := []struct {
		name  string
		cards Cards
		seq   Sequence
	}{
		{
			"no cards",
			Cards{},
			Sequence{Invalid, 0},
		},
		{
			"one card is ambiguous",
			Cards{1},
			Sequence{Ambiguous, 1},
		},
		{
			"two card straight",
			Cards{5, 6},
			Sequence{Straight, 2},
		},
		{
			"two card kind",
			Cards{5, 5},
			Sequence{Kind, 2},
		},
		{
			"five card valid straight",
			Cards{13, 13, 5, 6, 13},
			Sequence{Straight, 5},
		},
		{
			"five card invalid straight",
			Cards{13, 13, 2, 3, 13},
			Sequence{Invalid, 5},
		},
		{
			"five card valid straight",
			Cards{1, 2, 3, 4, 13},
			Sequence{Straight, 5},
		},
		{
			"invalid straight too high",
			Cards{11, 12, 13},
			Sequence{Invalid, 3},
		},
		{
			"3 of a kind with joker",
			Cards{13, 7, 7},
			Sequence{Kind, 3},
		},
		{
			"5 of a kind with joker",
			Cards{2, 2, 13, 2, 2},
			Sequence{Kind, 5},
		},
		{
			"ambiguous through jokers",
			Cards{13, 13, 5, 13, 13},
			Sequence{Ambiguous, 5},
		},
		{
			"all jokers is ambiguous",
			Cards{13, 13, 13, 13},
			Sequence{Ambiguous, 4},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			seq := Validate(c.cards)

			if seq.Type != c.seq.Type {
				t.Errorf("expected type %v but got %v", c.seq.Type, seq.Type)
			}
			if seq.N != c.seq.N {
				t.Errorf("expected N=%v but got N=%v", c.seq.N, seq.N)
			}
		})
	}
}
