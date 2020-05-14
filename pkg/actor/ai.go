package actor

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/nilsbu/phase10/pkg/game"
)

type AI struct{}

func (h *AI) Play(g *game.Game) error {
	if err := h.draw(g); err != nil {
		return err
	}

	if err := h.put(g); err != nil {
		return err
	}

	if len(g.Players[g.Turn].Cards) == 0 {
		g.Turn = (g.Turn + 1) % len(g.Players)
		return nil
	}

	if err := h.drop(g); err != nil {
		return err
	}

	return nil
}

func (h *AI) draw(g *game.Game) error {
	var fromTrash bool
	if g.Players[g.Turn].Out {
		fromTrash = false
	} else {
		fromTrash = h.shouldDrawTrashOut(g)
	}

	// if fromTrash {
	// 	fmt.Printf("%v drew %v from trash\n",
	// 		g.Players[g.Turn].Name, display.PrintCard(g.Trash, true))
	// } else {
	// 	fmt.Printf("%v drew from stack and left %v on the stack\n",
	// 		g.Players[g.Turn].Name, display.PrintCard(g.Trash, true))
	// }
	oldT := g.Trash
	scoresX := scoreCards(g.Players[g.Turn].Cards, g.Players[g.Turn].Phase)
	g.Draw(fromTrash)

	if fromTrash {
		scores := scoreCards(g.Players[g.Turn].Cards, g.Players[g.Turn].Phase)
		minScore := 1e+20
		minI := 0
		for i := range scores {
			if minScore > scores[i] {
				minScore = scores[i]
				minI = i
			}
		}
		if g.Players[g.Turn].Cards[minI] == oldT {
			fmt.Println("!!!!!!!!!!!!!!!!!!")
			// fmt.Println("1-13:", xxx, vt, exp)
			fmt.Println("cards:", g.Players[g.Turn].Cards)
			fmt.Println("scpost:", scores, minScore, minI)
			fmt.Println("scpre:", scoresX)
		}
	}

	return nil
}

func (h *AI) shouldDrawTrashOut(g *game.Game) bool {
	seqs, _ := game.GetPhaseSequences(g.Players[g.Turn].Phase)
	if _, ok := comeOutR(g.Players[g.Turn].Cards, seqs); ok {
		return isAppendable(g.Trash, g.OutCards)
	}

	cards := game.Cards{}
	for _, card := range g.Players[g.Turn].Cards {
		cards = append(cards, card)
	}
	cards = append(cards, g.Trash)
	if _, ok := comeOutR(cards, seqs); ok {
		return true
	}

	exp := 0.0
	vt := 0.0
	xxx := []float64{}
	for i := 1; i <= 13; i++ {
		cards := game.Cards{}
		for _, card := range g.Players[g.Turn].Cards {
			cards = append(cards, card)
		}
		cards = append(cards, game.Card(i))
		sort.Sort(cards)

		scores := scoreCards(cards, g.Players[g.Turn].Phase)
		sum := 0.0
		for _, s := range scores {
			sum += s
		}
		exp += sum / 13.0
		xxx = append(xxx, sum)

		if g.Trash == game.Card(i) {
			vt = sum
		}
	}

	return vt >= exp
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
		fmt.Printf("%v came out\n", g.Players[g.Turn].Name)
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

func isAppendable(card game.Card, out []game.Cards) bool {
	for _, seq := range out {
		seqCopy := game.Cards{card}
		for _, c := range seq {
			seqCopy = append(seqCopy, c)
		}

		if seq := game.Validate(seqCopy); seq.Type != game.Invalid {
			return true
		}

		seqCopy = append(seqCopy[1:], card)

		if seq := game.Validate(seqCopy); seq.Type != game.Invalid {
			return true
		}
	}

	return false
}

func (h *AI) drop(g *game.Game) error {
	scores := scoreCards(g.Players[g.Turn].Cards, g.Players[g.Turn].Phase)
	prefCards := sortByScore(g.Players[g.Turn].Cards, scores)

	if g.Players[(g.Turn+1)%len(g.Players)].Out {
		for _, card := range prefCards {
			if !isAppendable(card, g.OutCards) {
				return g.Drop(card)
			}
		}
	}

	return g.Drop(prefCards[0])
}

type cardScore struct {
	card  game.Card
	score float64
}

type cardScores []cardScore

func (cs cardScores) Len() int           { return len(cs) }
func (cs cardScores) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs cardScores) Less(i, j int) bool { return cs[i].score < cs[j].score }

func sortByScore(cards game.Cards, scores []float64) game.Cards {
	cs := cardScores{}
	for i := range cards {
		cs = append(cs, cardScore{cards[i], scores[i]})
	}

	sort.Sort(cs)

	out := game.Cards{}
	for _, pair := range cs {
		out = append(out, pair.card)
	}

	return out
}

func scoreCards(cards game.Cards, phase int) []float64 {
	cs := countCards(cards)
	sets := getFinishedSets(phase)

	iJokers := cs[game.CardTypes-1]

	scores := make([]float64, len(cards))
	for _, set := range sets {
		needed := setDiff(set, cs)
		chance := initialCollectChance(needed, iJokers)

		cNum := 0
		last := game.Card(0)
		for i := range cards {
			if cards[i] == 13 {
				scores[i] += chance
				continue
			}
			if last == cards[i] {
				cNum++
			} else {
				cNum = 1
			}
			last = cards[i]
			if set[int(cards[i])-1] >= cNum {
				scores[i] += chance
			}
		}
	}

	return scores
}

func setDiff(wanted, has []int) []int {
	diff := make([]int, len(wanted))
	for i := range wanted {
		if wanted[i] >= has[i] {
			diff[i] = wanted[i] - has[i]
		}
	}
	return diff
}

func countCards(cards game.Cards) []int {
	cs := make([]int, game.CardTypes)
	for _, c := range cards {
		cs[c-1]++
	}
	return cs
}

func getFinishedSets(phase int) (sets [][]int) {
	seqs, _ := game.GetPhaseSequences(phase)

	loneSets := make([][][]int, len(seqs))
	n := 1
	for s, seq := range seqs {
		switch seq.Type {
		case game.Kind:
			n *= 12
			for i := 0; i < 12; i++ {
				set := make([]int, 12)
				set[i] = seq.N
				loneSets[s] = append(loneSets[s], set)
			}
		case game.Straight:
			n *= 12 - seq.N + 1
			for i := 0; i < 12-seq.N+1; i++ {
				set := make([]int, 12)
				for j := 0; j < seq.N; j++ {
					set[i+j] = 1
				}
				loneSets[s] = append(loneSets[s], set)
			}
		}
	}

	is := make([]int, len(seqs))
	for j := 0; j < n; j++ {
		set := make([]int, 12)
		for s := range seqs {
			xset := loneSets[s][is[s]]
			for k := 0; k < 12; k++ {
				set[k] += xset[k]
			}
		}
		sets = append(sets, set)
		inc(is, loneSets)
	}

	return
}

func inc(is []int, sets [][][]int) {
	for j := len(is) - 1; j >= 0; j-- {
		is[j] = (is[j] + 1) % len(sets[j])
		if is[j] > 0 {
			return
		}
	}
}

func initialCollectChance(count []int, initialJokers int) float64 {
	sum := sum(count)

	c := 0
	for jokers := initialJokers; jokers <= sum; jokers++ {
		jCount := getJokerReplacements(count, jokers)
		for _, jc := range jCount {
			c += multiBinomial(jc)
		}
	}
	cProp := float64(c) * math.Pow(1.0/game.CardTypes, float64(sum-initialJokers))
	return cProp
}

func sum(count []int) int {
	s := 0
	for _, c := range count {
		s += c
	}
	return s
}

func multiBinomial(count []int) int {
	if len(count) == 0 {
		return 0
	}

	sum := sum(count)
	prod := 1
	for i := 0; i < len(count)-1; i++ {
		prod *= binomial(sum, count[i])
		sum -= count[i]
	}
	return prod
}

func binomial(n, k int) int {
	return int(big.NewInt(0).Binomial(int64(n), int64(k)).Int64())
}

func getJokerReplacements(count []int, jokers int) [][]int {
	if len(count) == 0 {
		return [][]int{}
	}

	jCount, _ := getJokerReplacementsR(count, jokers)
	for i := range jCount {
		jCount[i] = append(jCount[i], jokers)
	}
	return jCount
}

func getJokerReplacementsR(count []int, jokers int) (jCount [][]int, ok bool) {
	if len(count) == 0 {
		return [][]int{{}}, jokers == 0
	}

	for i := 0; i <= count[0]; i++ {
		jc, k := getJokerReplacementsR(count[1:], jokers-i)
		if !k {
			continue
		}

		for j := range jc {
			jc[j] = append([]int{count[0] - i}, jc[j]...)
		}
		jCount = append(jCount, jc...)
		ok = true
	}

	return
}
