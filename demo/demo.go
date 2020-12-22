package demo

import (
	"github.com/jessdwitch/spiders/battle"
	"github.com/jessdwitch/spiders/deck"
	"github.com/jessdwitch/spiders/engine"
	"github.com/jessdwitch/spiders/title"

	"github.com/hajimehoshi/ebiten/v2"
)

func BattleDemo() *engine.Game {
	const (
		screenWidth  = 640
		screenHeight = 480
	)

	state := engine.NewGameState(&title.TitleScene{})
	background := ebiten.NewImage(screenWidth, screenHeight)
	scene, err := battle.NewBattleScene(
		state,
		background,
		[]int{0, 0},
		true,
		map[deck.CardID]int{0: 5},
	)
	if err != nil {
		panic(err)
	}

	state.SceneManager.GoTo(scene)

	g := &engine.Game{nil, state, nil}

	return g
}
