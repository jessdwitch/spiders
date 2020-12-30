package demo

import (
	"encoding/csv"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jessdwitch/spiders/engine"
	"github.com/jessdwitch/spiders/engine/render"
)

type renderDemoScene struct {
	sprites []render.Sprite
}

func RenderDemo() *engine.Game {
	const (
		screenWidth  = 640
		screenHeight = 480
	)
	background := ebiten.NewImage(screenWidth, screenHeight)
	background.Fill(color.Black)

	sheetF, err := os.Open("./content/sprite/sheets.csv")
	if err != nil {
		panic(err)
	}
	defer sheetF.Close()
	sheetMani := csv.NewReader(sheetF)
	spriteF, err := os.Open("./content/sprite/sprites.csv")
	if err != nil {
		panic(err)
	}
	defer spriteF.Close()
	spriteMani := csv.NewReader(spriteF)

	sprites, err := render.NewSpriteFactoryFromManifests(sheetMani, spriteMani)
	if err != nil {
		panic(err)
	}
	scene, err := newRenderDemoScene(sprites)
	if err != nil {
		panic(err)
	}
	state := engine.NewGameState(scene)

	g := &engine.Game{nil, state, sprites}

	return g
}

func newRenderDemoScene(factory *render.SpriteFactory) (*renderDemoScene, error) {
	var err error
	var s render.Sprite
	sprites := make([]render.Sprite, 4)
	s, err = factory.GetSprite("slime_blue")
	if err != nil {
		return nil, err
	}
	s.Scale(4, 4)
	s.Translate(200, 200)
	s.Animate("idle")
	sprites[0] = s

	sprites[1], err = factory.GetSprite("slime_red")
	if err != nil {
		return nil, err
	}
	sprites[1].Scale(3, 3)
	sprites[1].Translate(300, 300)
	sprites[1].Animate("idle")
	sprites[1].SetDelay(10)
	sprites[2], err = factory.GetSprite("slime_green")
	if err != nil {
		return nil, err
	}
	sprites[2].Scale(2, 2)
	sprites[2].Translate(400, 200)
	sprites[2].Animate("idle")
	sprites[2].SetDelay(3)
	sprites[3], err = factory.GetSprite("slime_white")
	if err != nil {
		return nil, err
	}
	sprites[3].Translate(500, 300)
	sprites[3].Animate("idle")
	return &renderDemoScene{sprites}, nil
}

func (r *renderDemoScene) Draw(i *ebiten.Image) {
	for _, s := range r.sprites {
		s.Draw(i)
	}
}

func (r *renderDemoScene) Update(_ *engine.GameState) error {
	var err error
	for _, s := range r.sprites {
		if err = s.Update(); err != nil {
			return err
		}
	}
	return nil
}
