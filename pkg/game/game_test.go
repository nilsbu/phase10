package game

import (
	"math/rand"
	"reflect"
	"sort"
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

	err := game.Drop(12)
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
		cardss     CardSequence
		cardsAfter Cards
		outAfter   CardSequence
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
			[]Cards{{5, 5, 5}, {8, 13, 8, 8}},
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
			[]Cards{{1, 1, 1}, {11, 11, 13}},
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
			[]Cards{{1, 1, 1, 1}, {5, 5, 5}},
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
			[]Cards{{9, 9, 9}, {11, 11, 11}},
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
			[]Cards{{1, 1, 1, 1}, {4, 5, 7}},
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
						Cards: Cards{1, 1, 1, 1, 5, 5, 6, 7, 7, 10, 10},
						Phase: 1, Out: false}},
				OutCards: []Cards{{9, 9, 9}, {11, 11, 11}},
				Turn:     1, Trash: -1,
			},
			[]Cards{{1, 1, 1, 1}, {5, 5, 5}},
			Cards{1, 1, 1, 1, 5, 5, 6, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"already out",
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
			[]Cards{{1, 1, 1}, {5, 5, 5}},
			Cards{1, 1, 1, 1, 5, 5, 5, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"phase 1 not fulfilled",
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
			[]Cards{{1, 1, 1}, {5, 5, 5}},
			Cards{1, 1, 1, 1, 4, 5, 6, 7, 7, 10, 10},
			[]Cards{{9, 9, 9}, {11, 11, 11}},
			true,
		},
		{
			"phase 1 mixed order",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{2, 3, 5, 5, 7, 8, 9, 10, 13, 13, 13},
						Phase: 1, Out: false}},
				OutCards: []Cards{},
				Turn:     1, Trash: -1,
			},
			[]Cards{{3, 13, 13}, {5, 5, 13}},
			Cards{2, 7, 8, 9, 10},
			[]Cards{{3, 13, 13}, {5, 5, 13}},
			false,
		},
		{
			"phase 2: cards ambiguous",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 7, 8},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{2, 2, 4, 5, 9, 11, 11, 12, 13, 13, 13},
						Phase: 2, Out: false}},
				OutCards: []Cards{},
				Turn:     1, Trash: -1,
			},
			[]Cards{{2, 13, 4, 5}, {11, 11, 13, 13}},
			Cards{2, 9, 12},
			[]Cards{{2, 13, 4, 5}, {11, 11, 13, 13}},
			false,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			player := c.game.Turn
			err := c.game.ComeOut(c.cardss)
			if err != nil && !c.err {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err {
				t.Errorf("expected error but none occurred")
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
		card       Card
		seq        int
		left       bool
		cardsAfter Cards
		outAfter   CardSequence
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
			9, 1, false,
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
			5, 1, true,
			Cards{9, 10, 12},
			[]Cards{{1, 2, 3}, {5, 6, 13, 8}},
			false,
		},
		{
			"no joker in cards",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			13, 1, true,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"no 3 in cards",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: true}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			3, 1, true,
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
			5, 2, true,
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
			5, -1, true,
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
			5, 1, false,
			Cards{5, 9, 10, 12},
			[]Cards{{1, 2, 3}, {6, 13, 8}},
			true,
		},
		{
			"player is not out",
			&Game{
				Players: []Player{p1, {Name: "P2",
					Cards: Cards{5, 9, 10, 12},
					Phase: 10, Out: false}},
				OutCards: []Cards{{1, 2, 3}, {6, 13, 8}},
				Turn:     1, Trash: 11,
			},
			5, 1, true,
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
			"done when one player has no cards",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{},
						Phase: 9, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, -1,
		},
		{
			"done when all are out",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{3, 3},
						Phase: 9, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, -1,
		},
		{
			"done when one player has no cards with winner",
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
		{
			"done when all are out with winner",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 1, Out: true},
					{Name: "P2",
						Cards: Cards{3, 3},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, 1,
		},
		{
			"winner ambiguous",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 10, Out: true},
					{Name: "P2",
						Cards: Cards{3, 3},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, -1,
		},
		{
			"done when one player has no cards with winner when all are phase 10",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 4, 5, 5, 5, 6, 8, 8, 8, 13},
						Phase: 10, Out: true},
					{Name: "P2",
						Cards: Cards{},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, 1,
		},
		{
			"all done all phase 10",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{},
						Phase: 10, Out: true},
					{Name: "P2",
						Cards: Cards{},
						Phase: 10, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, -1,
		},
		{
			"all done all phase 10",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{6, 7, 8},
						Phase: 10, Out: false},
					{Name: "P2",
						Cards: Cards{},
						Phase: 8, Out: true}},
				OutCards: []Cards{},
				Turn:     0, Trash: 11,
			},
			true, -1,
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			done, winner := c.game.IsDone(), c.game.GetWinner()
			if done {
				if !c.done {
					t.Fatalf("shouldn't be done but is")
				}
			} else {
				if c.done {
					t.Fatalf("should be done but isn't")
				}
			}
			if winner != c.winner {
				t.Errorf("expected winner to be %v but is %v", c.winner, winner)
			}
		})
	}
}

func TestNextRound(t *testing.T) {
	cs := []struct {
		name     string
		gamePre  *Game
		gamePost *Game
	}{
		{
			"",
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{1, 2, 3, 4, 5, 6, 7, 8, 9, 9},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{10, 11, 13},
						Phase: 3, Out: true}},
				OutCards: []Cards{{1, 1, 1}},
				Turn:     1, Trash: 11,
			},
			&Game{
				Players: []Player{
					{Name: "P1",
						Cards: Cards{2, 3, 3, 4, 4, 6, 7, 8, 12, 13},
						Phase: 1, Out: false},
					{Name: "P2",
						Cards: Cards{2, 3, 4, 4, 7, 7, 7, 8, 11, 13},
						Phase: 4, Out: false}},
				OutCards: []Cards{},
				Turn:     1, Trash: 6,
			},
		},
	}

	for _, c := range cs {
		t.Run(c.name, func(t *testing.T) {
			rand.Seed(0)
			c.gamePre.NextRound()

			if !reflect.DeepEqual(c.gamePre, c.gamePost) {
				t.Errorf("game is wrong\nwant: %v\nhas:  %v",
					c.gamePost, c.gamePre)
			}
		})
	}
}

const J = Card(13)

func TestSortSequences(t *testing.T) {
	cs := []struct {
		seq, seqSorted CardSequence
	}{
		{CardSequence{{J, 1, J}},
			CardSequence{{J, 1, J}}},
		{CardSequence{{3, 3, 3}, {J, 1, J}},
			CardSequence{{J, 1, J}, {3, 3, 3}}},
		{CardSequence{{1, 1, 1, 1}, {J, 1, J}},
			CardSequence{{1, 1, 1, 1}, {J, 1, J}}},
		{CardSequence{{1, 1, 1, 2}, {J, 1, J}},
			CardSequence{{J, 1, J}, {1, 1, 1, 2}}},
		{CardSequence{},
			CardSequence{}},
		{CardSequence{{3, J, J}, {3, 4, 5, 6}},
			CardSequence{{3, J, J}, {3, 4, 5, 6}}},
		{CardSequence{{3, J, 4}, {3, 3, 5, 6}},
			CardSequence{{3, 3, 5, 6}, {3, J, 4}}},
		{CardSequence{{3, 4, 5, 6}, {3, 4, 5, 6, 7}},
			CardSequence{{3, 4, 5, 6, 7}, {3, 4, 5, 6}}},
		{CardSequence{{3, 4, 5, 6}, {4, 5, 6, 7}, {1, 2, 3}},
			CardSequence{{1, 2, 3}, {3, 4, 5, 6}, {4, 5, 6, 7}}},
		{CardSequence{{3, 4, 5, 6}, {4, 5, 6, 7}, {}},
			CardSequence{{3, 4, 5, 6}, {4, 5, 6, 7}, {}}},
		{CardSequence{{13, 4, 5, 6}, {4, 5, 6, 7}},
			CardSequence{{13, 4, 5, 6}, {4, 5, 6, 7}}},
		{CardSequence{{13, 13, 4}, {4, 5, 6, 7}},
			CardSequence{{13, 13, 4}, {4, 5, 6, 7}}},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			sort.Sort(c.seq)

			if !reflect.DeepEqual(c.seq, c.seqSorted) {
				t.Errorf("game is wrong\nwant: %v\nhas:  %v",
					c.seq, c.seqSorted)
			}
		})
	}
}
