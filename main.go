package main

import (
	"log"

	"github.com/jessdwitch/spiders/demo"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tacocat")
	if err := ebiten.RunGame(demo.BattleDemo()); err != nil {
		log.Fatal(err)
	}
}
