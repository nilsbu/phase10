package game

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestServeN(t *testing.T) {
	rand.Seed(0)
	{
		cards := serveN(10)
		expected := Cards{2, 3, 3, 4, 4, 6, 7, 8, 12, 13}
		if !reflect.DeepEqual(cards, expected) {
			t.Errorf("expected cards %v but got %v", expected, cards)
		}
	}
	{
		cards := serveN(10)
		expected := Cards{2, 3, 4, 4, 7, 7, 7, 8, 11, 13}
		if !reflect.DeepEqual(cards, expected) {
			t.Errorf("expected cards %v but got %v", expected, cards)
		}
	}
}

func TestSetUp(t *testing.T) {
	rand.Seed(0)
	game := SetUp(2, 1)

	{
		if game.Players[0].Name != "Player 1" {
			t.Errorf("expected the 1st player's name to be 'Player 1' but got '%v'",
				game.Players[0].Name)
		}

		expected := Cards{2, 3, 3, 4, 4, 6, 7, 8, 12, 13}
		if !reflect.DeepEqual(game.Players[0].Cards, expected) {
			t.Errorf("expected the 1st player's cards to be '%v' but got '%v'",
				expected, game.Players[0].Cards)
		}

		if game.Players[0].Phase != 1 {
			t.Errorf("expected the 1st player to be in phase 1 but got '%v'",
				game.Players[0].Phase)
		}

		if game.Players[0].Out != false {
			t.Error("player 1 must not be out")
		}
	}
	{
		if game.Players[1].Name != "Player 2" {
			t.Errorf("expected the 2nd player's name to be 'Player 1' but got '%v'",
				game.Players[1].Name)
		}

		expected := Cards{2, 3, 4, 4, 7, 7, 7, 8, 11, 13}
		if !reflect.DeepEqual(game.Players[1].Cards, expected) {
			t.Errorf("expected the 2nd player's cards to be '%v' but got '%v'",
				expected, game.Players[1].Cards)
		}

		if game.Players[1].Phase != 1 {
			t.Errorf("expected the 2nd player to be in phase 1 but got '%v'",
				game.Players[1].Phase)
		}

		if game.Players[1].Out != false {
			t.Error("player 2 must not be out")
		}
	}

	if game.Turn != 1 {
		t.Errorf("expected it to be 1's turn but it's %v's",
			game.Turn)
	}

	if game.Trash != 6 {
		t.Errorf("expected 6 in the trash but got %v",
			game.Trash)
	}
}

func TestDrawFromTrash(t *testing.T) {
	rand.Seed(0)
	game := SetUp(2, 0)
	game.Draw(true)

	expected := Cards{2, 3, 3, 4, 4, 6, 6, 7, 8, 12, 13} // 6 got added
	if !reflect.DeepEqual(game.Players[0].Cards, expected) {
		t.Errorf("expected the 1st player's cards to be '%v' but got '%v'",
			expected, game.Players[0].Cards)
	}
	if game.Trash != -1 {
		t.Errorf("expected -1 in the trash but got %v",
			game.Trash)
	}
}

func TestDrawFromStackAndDrop(t *testing.T) {
	rand.Seed(0)
	game := SetUp(2, 0)
	game.Draw(false)

	expected := Cards{2, 3, 3, 4, 4, 6, 6, 7, 8, 12, 13} // 6 got added
	if !reflect.DeepEqual(game.Players[0].Cards, expected) {
		t.Errorf("expected the 1st player's cards to be '%v' but got '%v'",
			expected, game.Players[0].Cards)
	}
	if game.Trash != 6 {
		t.Errorf("expected 6 in the trash but got %v",
			game.Trash)
	}
	if game.Turn != 0 {
		t.Errorf("expected it to be 1's turn but it's %v's",
			game.Turn)
	}

	err := game.Drop(9)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	expected = Cards{2, 3, 3, 4, 4, 6, 6, 7, 8, 13} // 12 got dropped
	if !reflect.DeepEqual(game.Players[0].Cards, expected) {
		t.Errorf("expected the 1st player's cards to be '%v' but got '%v'",
			expected, game.Players[0].Cards)
	}
	if game.Trash != 12 {
		t.Errorf("expected 12 in the trash but got %v",
			game.Trash)
	}
	if game.Turn != 1 {
		t.Errorf("expected it to be 1's turn but it's %v's",
			game.Turn)
	}

	err = game.Drop(-1)
	if err == nil {
		t.Fatalf("expected error for negative index but none occurred")
	}
}

func TestComeOut(t *testing.T) {
	p1 := Player{Name: "P1",
		Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
		Phase: 1, Out: false}

	cs := []struct {
		name       string
		game       *Game
		idxSeq     [][]int
		cardsAfter Cards
		outAfter   []Cards
		err        bool
	}{
		{
			"phase 1",
			&Game{
				Players: []Player{
					p1,
					{Name: "P2",
						Cards: Cards{10, 11, 13},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			[][]int{{3, 4, 5}, {7, 10, 8, 9}},
			Cards{2, 3, 4, 6},
			[]Cards{{5, 5, 5}, {8, 13, 8, 8}},
			false,
		},
		{
			"phase 2", // TODO is that phase 2?
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{1, 1, 1, 2, 6, 9, 9, 10, 11, 11, 13},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{10, 11, 13},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			[][]int{{0, 1, 2}, {8, 9, 10}},
			Cards{2, 6, 9, 9, 10},
			[]Cards{{1, 1, 1}, {11, 11, 13}},
			false,
		},
		{
			"phase 1+",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
						Phase: 1, Out: false}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[][]int{{0, 1, 2, 3}, {4, 5, 6}},
			Cards{7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}, {1, 1, 1, 1}, {5, 5, 5}},
			false,
		},
		{
			"no  cards",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
						Phase: 1, Out: false}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[][]int{},
			Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"invalid sequence",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{1, 1, 1, 1, 4, 5, 6, 7, 7, 10, 10},
						Phase: 1, Out: false}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[][]int{{0, 1, 2, 3}, {4, 5, 7}},
			Cards{1, 1, 1, 1, 4, 5, 6, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"duplicate card",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
						Phase: 1, Out: false}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[][]int{{0, 1, 2, 3}, {4, 5, 5}},
			Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"phase 1+",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
						Phase: 1, Out: true}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[][]int{{0, 1, 2, 3}, {4, 5, 6}},
			Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			player := c.game.Turn
			err := c.game.ComeOut(c.idxSeq)
			if (err != nil) != c.err {
				t.Errorf("wrong error")
			}
			if !reflect.DeepEqual(c.game.Players[player].Cards, c.cardsAfter) {
				t.Errorf("expected the player's cards to be '%v' but got '%v'",
					c.cardsAfter, c.game.Players[player].Cards)
			}
			if !reflect.DeepEqual(c.game.OutCards, c.outAfter) {
				t.Errorf("expected out cards to be '%v' but got '%v'",
					c.outAfter, c.game.OutCards)
			}
			if err == nil && !c.game.Players[player].Out {
				t.Errorf("player is not out")
			}
		})
	}
}

func TestAppend(t *testing.T) {
	p1 := Player{Name: "P1",
		Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
		Phase: 1, Out: false}

	cs := []struct {
		name       string
		game       *Game
		card, seq  int
		left       bool
		cardsAfter Cards
		outAfter   []Cards
		err        bool
	}{
		{
			"append right",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			1, 1, false,
			Cards{5, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8, 9}},
			false,
		},
		{
			"append left",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			0, 1, true,
			Cards{9, 10, 12},
			[]Cards{{1, 2, 3}, {5, 6, 13, 8}},
			false,
		},
		{
			"card oob right",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			10, 1, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"card oob left",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			-1, 1, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"sequence oob right",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			0, 2, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"sequence oob left",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			0, -1, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"invalid append",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			0, 1, false,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"append left",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: false}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			0, 1, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			player := c.game.Turn
			err := c.game.Append(c.card, c.seq, c.left)
			if (err != nil) != c.err {
				t.Errorf("wrong error")
			}

			if !reflect.DeepEqual(c.game.Players[player].Cards, c.cardsAfter) {
				t.Errorf("expected the player's cards to be '%v' but got '%v'",
					c.cardsAfter, c.game.Players[player].Cards)
			}
			if !reflect.DeepEqual(c.game.OutCards, c.outAfter) {
				t.Errorf("expected out cards to be '%v' but got '%v'",
					c.outAfter, c.game.OutCards)
			}
		})
	}
}

func TestIsDone(t *testing.T) {
	cs := []struct {
		name   string
		game   *Game
		done   bool
		winner int
	}{
		{
			"not done",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{10, 11, 13},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			false, -1,
		},
		{
			"done",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, 1,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			done, winner := c.game.IsDone(), c.game.GetWinner()
			if done {
				if !c.done {
					t.Fatalf("shouldn't be done but is")
				}
				if winner != c.winner {
					t.Errorf("expected winner to be %v but is %v", c.winner, winner)
				}
			} else {
				if c.done {
					t.Fatalf("should be done but isn't")
				}
			}
		})
	}
}
