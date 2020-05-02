package actor

import (
	"bufio"
	"errors"
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
	if err := h.draw(g); err != nil {
		return err
	}
	if err := h.put(g); err != nil {
		return err
	}
	if err := h.drop(g); err != nil {
		return err
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
		for _, cmd := range cmds {
			g.Append(cmd.cardIdx-1, cmd.sequenceIdx-1, cmd.left)
		}

	} else {
		text, err := read("come out", `^((([\d]+,)*[\d]+;)+)|-$`)
		if err != nil {
			return err
		} else if text == "-" {
			return nil
		}

		cmds := splitCOString(text)
		for _, cmd := range cmds {
			for i := 0; i < len(cmd); i++ {
				cmd[i]--
			}
		}
		g.ComeOut(cmds)

		h.put(g)
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
	idx, _ := strconv.Atoi(text)
	g.Drop(idx - 1)
	return nil
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

		if text == "quit" {
			return "", errors.New("quit")
		} else if re.MatchString(text) {
			return text, nil
		}
	}
}
