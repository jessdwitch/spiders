package title

import (
	"github.com/jessdwitch/spiders/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type TitleScene struct {
	count int
}

func anyGamepadAbstractButtonPressed(i *engine.Input) bool {
	return false // TODO
}

func (s *TitleScene) Update(state *engine.GameState) error {
	return nil // TODO
}

func (s *TitleScene) Draw(r *ebiten.Image) {

}
