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
		SceneManager: NewSceneManager(initScene),
		Input:        NewInput(),
	}
}

func (g *Game) Update() error {
	return nil // TODO
}

func (g *Game) Draw(*ebiten.Image) {
	// TODO
}

func (g *Game) Layout(outisdeWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 0, 0 // TODO
}
