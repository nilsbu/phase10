package actor

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/nilsbu/phase10/pkg/display"
	"github.com/nilsbu/phase10/pkg/game"
)

type Human struct{}

func (h *Human) Play(g *game.Game) error {
	for {
		if err := h.draw(g); err == nil {
			break
		} else {
			fmt.Println(err)
		}
	}
	for {
		if err := h.put(g); err == nil {
			break
		} else {
			fmt.Println(err)
		}
	}
	for {
		if err := h.drop(g); err == nil {
			break
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func (h *Human) draw(g *game.Game) error {
	fmt.Println(display.PrintGame(g, g.Turn))
	text, err := read("draw card (s/t)", `^[st]$`)
	if err != nil {
		return err
	}

	g.Draw(text == "t")
	return nil
}

func (h *Human) put(g *game.Game) error {
	fmt.Println(display.PrintGame(g, g.Turn))
	if g.Players[g.Turn].Out {
		text, err := read("append", `^(([\d]+[><][\d]+;)+)|-$`)
		if err != nil {
			return err
		} else if text == "-" {
			return nil
		}

		cmds := splitAddString(text)
		cards := []int{}
		for _, cmd := range cmds {
			cards = append(cards, cmd.cardIdx)
		}
		cardIdxs, err := findCards(cards, g.Players[g.Turn].Cards)
		if err != nil {
			return err
		}
		for i, cmd := range cmds {
			if err = g.Append(cardIdxs[i], cmd.sequenceIdx-1, cmd.left); err != nil {
				return err
			}
			for j := i + 1; j < len(cmds); j++ {
				if cardIdxs[j] > cardIdxs[i] {
					cardIdxs[j]--
				}
			}
		}

	} else {
		text, err := read("come out", `^((([\d]+,)*[\d]+;)+)|-$`)
		if err != nil {
			return err
		} else if text == "-" {
			return nil
		}

		cmds := splitCOString(text)
		cardIdxs := [][]int{}
		for _, cmd := range cmds {
			idxs, err := findCards(cmd, g.Players[g.Turn].Cards)
			if err != nil {
				return err
			}
			cardIdxs = append(cardIdxs, idxs)
		}

		if err := g.ComeOut(cardIdxs); err != nil {
			return err
		}
		if err := h.put(g); err != nil {
			return err
		}
	}
	return nil
}

type addCmd struct {
	cardIdx, sequenceIdx int
	left                 bool
}

func splitAddString(s string) (cmds []addCmd) {
	ss := strings.SplitN(s, ";", -1)
	re := regexp.MustCompile(`[<>]`)
	for _, str := range ss {
		if len(str) < 3 {
			continue
		}
		nums := re.Split(str, 2)

		card, _ := strconv.Atoi(nums[0])
		seq, _ := strconv.Atoi(nums[1])

		cmds = append(cmds, addCmd{card, seq, strings.Contains(str, "<")})
	}
	return
}

func splitCOString(s string) (cmds [][]int) {
	ss := strings.SplitN(s, ";", -1)
	for _, str := range ss {
		if len(str) < 2 {
			continue
		}

		numStr := strings.SplitN(str, ",", -1)
		var nums []int

		for _, ns := range numStr {
			num, _ := strconv.Atoi(ns)
			nums = append(nums, num)
		}

		cmds = append(cmds, nums)
	}
	return
}

func (h *Human) drop(g *game.Game) error {
	fmt.Println(display.PrintGame(g, g.Turn))
	text, err := read("drop card", `^[\d]+$`)
	if err != nil {
		return err
	}
	card, _ := strconv.Atoi(text)
	idxs, err := findCards([]int{card}, g.Players[g.Turn].Cards)
	if err != nil {
		return err
	}
	return g.Drop(idxs[0])
}

func read(prompt, regex string) (string, error) {
	re := regexp.MustCompile(regex)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%v: ", prompt)
		text, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		text = text[:len(text)-1]

		if re.MatchString(text) {
			return text, nil
		}
	}
}

func findCards(want []int, has game.Cards) (idxs []int, err error) {
	found := make([]bool, len(has))

	for _, w := range want {
		f := false
		for i, h := range has {
			if game.Card(w) == h && !found[i] {
				found[i] = true
				idxs = append(idxs, i)
				f = true
				break
			}
		}
		if !f {
			return []int{}, fmt.Errorf("card '%v' not found", w)
		}
	}

	return
}
