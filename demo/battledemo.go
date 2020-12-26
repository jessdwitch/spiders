package demo

import (
	"encoding/csv"
	"image/color"
	"os"

	"github.com/jessdwitch/spiders/battle"
	"github.com/jessdwitch/spiders/deck"
	"github.com/jessdwitch/spiders/engine"
	"github.com/jessdwitch/spiders/engine/render"
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
	background.Fill(color.RGBA{240, 177, 177, 1})
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

	sheetF, err := os.Open("../static/sprite/sheets.csv")
	if err != nil {
		panic(err)
	}
	defer sheetF.Close()
	sheetMani := csv.NewReader(sheetF)
	spriteF, err := os.Open("../static/sprite/sprites.csv")
	if err != nil {
		panic(err)
	}
	defer spriteF.Close()
	spriteMani := csv.NewReader(spriteF)

	sprites, err := render.NewSpriteFactoryFromManifests(sheetMani, spriteMani)
	if err != nil {
		panic(err)
	}

	g := &engine.Game{nil, state, sprites}

	return g
}
