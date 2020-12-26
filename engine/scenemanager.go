// Manage current scene and handle transitions between them

package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	transitionFrom = ebiten.NewImage(1280, 960) //(ScreenWidth, ScreenHeight)
	transitionTo   = ebiten.NewImage(1280, 960) //(ScreenWidth, ScreenHeight)
)

const transitionMaxCount = 20

type (
	Scene interface {
		Update(state *GameState) error
		Draw(screen *ebiten.Image)
	}
	SceneManager struct {
		current         Scene
		next            Scene
		transitionCount int
	}
)

func NewSceneManager(initialScene Scene) *SceneManager {
	s := &SceneManager{
		current: initialScene,
	}
	return s
}

// Update : Call the current Scene's Update, or do nothing if in transition
func (s *SceneManager) Update(state *GameState) error {
	if s.transitionCount == 0 {
		return s.current.Update(state)
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil

	return nil
}

// Draw : Draw the current scene or handle the scene transition if need be
func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	r.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(transitionTo, op)
}

// GoTo : Initiate a scene transition to the given Scene
func (s *SceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}
