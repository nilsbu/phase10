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
		{
			game.Cards{1, 7, 8, 10, 11, 12, 13, 13, 13, 13},
			[]game.Sequence{{Type: game.Straight, N: 9}},
			[]game.Cards{{13, 13, 13, 7, 8, 13, 10, 11, 12}}, true,
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

// func TestScore(t *testing.T) {
// 	var a, b, c, J game.Card = 10, 11, 12, 13
// 	fmt.Sprintln(a, b, c, J)
// 	phase := 5
//
// 	cards := game.Cards{1, 7, 8, 10, 11, 12, 13, 13, 13, 13}
// 	scores := scoreCards(cards, phase)
// 	for i, card := range cards {
// 		fmt.Printf("%v: %v\n", card, scores[i])
// 	}
// 	fmt.Println()
//
// 	v := 0.0
// 	baseLine := 0.0
// 	for _, s := range scoreCards(cards, phase) {
// 		baseLine += s
// 	}
//
// 	for i := 1; i <= 13; i++ {
// 		xcards := game.Cards{}
// 		for _, card := range cards {
// 			xcards = append(xcards, card)
// 		}
// 		xcards = append(cards, game.Card(i))
// 		sort.Sort(xcards)
// 		scores := scoreCards(xcards, phase)
// 		sum := 0.0
// 		for _, s := range scores {
// 			// if cards[i] != 13 {
// 			sum += s
// 			// }
// 		}
// 		sum -= baseLine
// 		v += sum
// 		fmt.Println(i, sum)
// 	}
// 	fmt.Println(v / 13)
// }

func TestShhouldDrawTrash(t *testing.T) {
	cs := []struct {
		game *game.Game
		draw bool
	}{
		{
			&game.Game{
				Players: []game.Player{
					{Name: "P1",
						Cards: game.Cards{2, 3, 4, 6, 6, 7, 8, 11, 12, 13},
						Phase: 5, Out: false}},
				OutCards: []game.Cards{},
				Turn:     0, Trash: 6,
			},
			false,
		},
		{
			&game.Game{
				Players: []game.Player{
					{Name: "P1",
						Cards: game.Cards{1, 7, 8, 10, 11, 12, 13, 13, 13, 13},
						Phase: 6, Out: false}},
				OutCards: []game.Cards{},
				Turn:     0, Trash: 1,
			},
			false,
		},
		{
			&game.Game{
				Players: []game.Player{
					{Name: "P1",
						Cards: game.Cards{1, 7, 8, 10, 11, 12, 13, 13, 13, 13},
						Phase: 6, Out: false}},
				OutCards: []game.Cards{},
				Turn:     0, Trash: 13,
			},
			true,
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			ai := &AI{}
			draw := ai.shouldDrawTrash(c.game)
			if c.draw != draw {
				t.Errorf("%v != %v", c.draw, draw)
			}
		})
	}
}

func TestJokerReplacements(t *testing.T) {
	cs := []struct {
		count  []int
		jokers int
		jCount [][]int
	}{
		{
			[]int{}, 0,
			[][]int{},
		},
		{
			[]int{1}, 0,
			[][]int{{1, 0}},
		},
		{
			[]int{1}, 1,
			[][]int{{0, 1}},
		},
		{
			[]int{2}, 1,
			[][]int{{1, 1}},
		},
		{
			[]int{1, 1}, 1,
			[][]int{{1, 0, 1}, {0, 1, 1}},
		},
		{
			[]int{2, 1}, 2,
			[][]int{{1, 0, 2}, {0, 1, 2}},
		},
		{
			[]int{2, 1}, 0,
			[][]int{{2, 1, 0}},
		},
		{
			[]int{2, 1}, 3,
			[][]int{{0, 0, 3}},
		},
		{
			[]int{2, 1}, 1,
			[][]int{{2, 0, 1}, {1, 1, 1}},
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			jCount := getJokerReplacements(c.count, c.jokers)
			if !reflect.DeepEqual(c.jCount, jCount) {
				t.Errorf("%v != %v", c.jCount, jCount)
			}
		})
	}
}

func TestMultiBinomial(t *testing.T) {
	cs := []struct {
		count []int
		mb    int
	}{
		{[]int{}, 0},
		{[]int{1}, 1},
		{[]int{1, 1}, 2},
		{[]int{1, 0}, 1},
		{[]int{0, 1, 1, 0}, 2},
		{[]int{2, 2, 2}, 90},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			mb := multiBinomial(c.count)
			if c.mb != mb {
				t.Errorf("%v != %v", c.mb, mb)
			}
		})
	}
}

func TestIsAppendable(t *testing.T) {
	cs := []struct {
		card       game.Card
		out        []game.Cards
		appendable bool
	}{
		{1, []game.Cards{}, false},
		{1, []game.Cards{{13, 1, 1}}, true},
		{2, []game.Cards{{13, 1, 1}}, false},
		{4, []game.Cards{{1, 2, 3}}, true},
		{4, []game.Cards{{5, 5, 5}, {1, 13, 13}}, true},
		{4, []game.Cards{{5, 5, 5}, {13, 13, 13}}, true},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			appendable := isAppendable(c.card, c.out)
			if appendable && !c.appendable {
				t.Errorf("card '%v' must not be appendable", c.card)
			} else if !appendable && c.appendable {
				t.Errorf("card '%v' must be appendable", c.card)
			}
		})
	}
}

func TestSortByScore(t *testing.T) {
	cs := []struct {
		cards  game.Cards
		scores []float64
		sorted game.Cards
	}{
		{
			game.Cards{},
			[]float64{},
			game.Cards{},
		},
		{
			game.Cards{1},
			[]float64{1.0},
			game.Cards{1},
		},
		{
			game.Cards{1, 2},
			[]float64{1.0, .8},
			game.Cards{2, 1},
		},
		{
			game.Cards{1, 1, 2, 3},
			[]float64{1.0, 0.1, .8, .4},
			game.Cards{1, 3, 2, 1},
		},
		{
			game.Cards{1, 1, 2, 3},
			[]float64{1.0, 0.4, .8, .4},
			game.Cards{1, 3, 2, 1},
		},
	}

	for _, c := range cs {
		t.Run("", func(t *testing.T) {
			sorted := sortByScore(c.cards, c.scores)
			if !reflect.DeepEqual(c.sorted, sorted) {
				t.Errorf("%v != %v", c.sorted, sorted)
			}
		})
	}
}
