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
	if len(g.Players[g.Turn].Cards) == 0 {
		g.Turn = (g.Turn + 1) % len(g.Players)
		return nil
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
	text, err := read("draw card (s/t)", `^[st]?$`)
	if err != nil {
		return err
	}

	g.Draw(text == "t")
	return nil
}

func (h *Human) put(g *game.Game) error {
	fmt.Println(display.PrintGame(g, g.Turn))
	if g.Players[g.Turn].Out {
		text, err := read("append", `^((([\dabcJj]+[><][\dabcJj]+;)+)|-|)$`)
		if err != nil {
			return err
		} else if text == "-" || text == "" {
			return nil
		}

		cmds := splitAddString(text)
		cards := game.Cards{}
		for _, cmd := range cmds {
			cards = append(cards, cmd.card)
		}

		for i, cmd := range cmds {
			if err = g.Append(cards[i], cmd.sequenceIdx-1, cmd.left); err != nil {
				return err
			}
		}

	} else {
		text, err := read("come out", `^(((([\dabcJj]+,)*[\dabcJj]+;)+)|-|)$`)
		if err != nil {
			return err
		} else if text == "-" || text == "" {
			return nil
		}

		cmds := splitCOString(text)
		if err := g.ComeOut(cmds); err != nil {
			return err
		}
		if err := h.put(g); err != nil {
			return err
		}
	}
	return nil
}

type addCmd struct {
	card        game.Card
	sequenceIdx int
	left        bool
}

func splitAddString(s string) (cmds []addCmd) {
	ss := strings.SplitN(s, ";", -1)
	re := regexp.MustCompile(`[<>]`)
	for _, str := range ss {
		if len(str) < 3 {
			continue
		}
		nums := re.Split(str, 2)

		card := readCard(nums[0])
		seq, _ := strconv.Atoi(nums[1])

		cmds = append(cmds, addCmd{card, seq, strings.Contains(str, "<")})
	}
	return
}

func readCard(s string) game.Card {
	switch s {
	case "a":
		return 10
	case "b":
		return 11
	case "c":
		return 12
	case "J", "j":
		return 13
	default:
		i, _ := strconv.Atoi(s)
		return game.Card(i)
	}
}

func splitCOString(s string) (cmds []game.Cards) {
	ss := strings.SplitN(s, ";", -1)
	for _, str := range ss {
		if len(str) < 2 {
			continue
		}

		numStr := strings.SplitN(str, ",", -1)
		var nums game.Cards

		for _, ns := range numStr {
			nums = append(nums, readCard(ns))
		}

		cmds = append(cmds, nums)
	}
	return
}

func (h *Human) drop(g *game.Game) error {
	fmt.Println(display.PrintGame(g, g.Turn))
	text, err := read("drop card", `^[\dabcJj]+$`)
	if err != nil {
		return err
	}

	return g.Drop(game.Card(readCard(text)))
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
