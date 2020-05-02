package actor

import "github.com/nilsbu/phase10/pkg/game"

type Actor interface {
	Play(g *game.Game) error
}
