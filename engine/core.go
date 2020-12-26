// Entry point into the game loop

package engine

import (
	"github.com/jessdwitch/spiders/engine/render"

	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type (
	Game struct {
		Events       *eventBus
		GameState    *GameState
		SpriteGetter render.SpriteGetter
	}

	GameState struct {
		Config       Config
		SceneManager *SceneManager
		Input        *Input
		PlayerParty  PlayerParty
	}
)

// NewGame : Generate a new Game object.
func NewGame(initScene Scene) (*Game, error) {
	return &Game{
		Events: &eventBus{},
		// TODO: Initial scene
		GameState: NewGameState(initScene),
	}, nil
}

func NewGameState(initScene Scene) *GameState {
	return &GameState{
		Config:       Config{ScreenHeight: 480, ScreenWidth: 640},
		SceneManager: NewSceneManager(initScene),
		Input:        NewInput(),
	}
}

func (g *Game) Update() error {
	return g.GameState.SceneManager.Update(g.GameState)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.GameState.SceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
