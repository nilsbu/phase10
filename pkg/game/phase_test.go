package game

import "testing"

func TestIsPhaseFulfilled(t *testing.T) {
	cs := []struct {
		name      string
		seqs      []Sequence
		phase     int
		fulfilled bool
	}{
		{
			"invalid phase 0",
			[]Sequence{},
			0, false,
		},
		{
			"invalid phase 9999",
			[]Sequence{{Straight, 9}},
			9999, false,
		},
		{
			"phase 1: fulfilled",
			[]Sequence{{Kind, 3}, {Kind, 3}},
			1, true,
		},
		{
			"phase 1: overfulfilled one seq",
			[]Sequence{{Kind, 3}, {Kind, 5}},
			1, true,
		},
		{
			"phase 1: wrong sequence",
			[]Sequence{{Kind, 3}, {Straight, 3}},
			1, false,
		},
		{
			"phase 1: sequence too short",
			[]Sequence{{Kind, 2}, {Kind, 3}},
			1, false,
		},
		{
			"phase 1: extra sequences forbidden",
			[]Sequence{{Kind, 3}, {Kind, 3}, {Kind, 3}},
			1, false,
		},
		{
			"phase 1: invalid sequence",
			[]Sequence{{Kind, 3}, {Invalid, 3}},
			1, false,
		},
		{
			"phase 1: ambiguous sequence",
			[]Sequence{{Kind, 3}, {Ambiguous, 3}},
			1, true,
		},
		{
			"phase 1: no sequences",
			[]Sequence{},
			1, false,
		},
		{
			"phase 2: fulfilled",
			[]Sequence{{Kind, 3}, {Straight, 4}},
			2, true,
		},
		{
			"phase 2: fulfilled with reverse order",
			[]Sequence{{Straight, 5}, {Kind, 3}},
			2, true,
		},
		{
			"phase 2: straight too short",
			[]Sequence{{Kind, 3}, {Straight, 3}},
			2, false,
		},
		{
			"phase 2: kind missing",
			[]Sequence{{Straight, 5}},
			2, false,
		},
		// The following just test if the minimum of a phase is fulfilled
		{
			"phase 3",
			[]Sequence{{Kind, 4}, {Straight, 4}},
			3, true,
		},
		{
			"phase 4",
			[]Sequence{{Straight, 7}},
			4, true,
		},
		{
			"phase 5",
			[]Sequence{{Straight, 8}},
			5, true,
		},
		{
			"phase 6",
			[]Sequence{{Straight, 9}},
			6, true,
		},
		{
			"phase 7",
			[]Sequence{{Kind, 4}, {Kind, 4}},
			7, true,
		},
		{
			"phase 8",
			[]Sequence{{Kind, 5}, {Kind, 2}},
			8, true,
		},
		{
			"phase 9",
			[]Sequence{{Kind, 5}, {Kind, 3}},
			9, true,
		},
		{
			"phase 10",
			[]Sequence{{Kind, 5}, {Straight, 5}},
			10, true,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			fulfilled := isPhaseFulfilled(c.seqs, c.phase)

			if fulfilled && !c.fulfilled {
				t.Errorf("expected phase %v to not be fulfilled but was", c.phase)
			} else if !fulfilled && c.fulfilled {
				t.Errorf("expected phase %v to be fulfilled but was not", c.phase)
			}
		})
	}
}
