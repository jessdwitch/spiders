package render_test

import (
	"bytes"
	"encoding/csv"
	"image"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jessdwitch/spiders/engine/render"
	"github.com/jessdwitch/spiders/engine/render/testdata"

	"github.com/stretchr/testify/assert"

	_ "image/png"
)

func TestExtractAnimation(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(testdata.Test_sprite_png))
	if err != nil {
		t.Fatal(err)
	}
	s := &render.SpriteSheet{
		SourceImage: ebiten.NewImageFromImage(img),
		SheetDim:    image.Pt(204, 54),
		TileDim:     image.Pt(51, 54),
		NTiles:      image.Pt(4, 1),
	}
	meta := sampleAnimMetas[0]
	anim, err := s.ExtractAnimation(meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, anim.CurrentFrame)
	assert.Equal(t, 0, anim.FrameDelay)
	assert.Equal(t, meta.FrameDelay, anim.FrameDelayMax)
	assert.Len(t, anim.Frames, meta.NFrames)
	assert.NotContains(t, anim.Frames, nil)
	assert.Equal(t, image.Pt(51, 54), anim.GetDims())
	assert.Equal(t, render.Point{0, 0}, anim.GetPosition())
}

// func TestSpriteSheetFiles(t *testing.T) {
// 	s, err := render.NewSpriteSheetFiles(csv.NewReader(strings.NewReader(sampleSheetManifest)))
// 	assert.NoError(t, err)
// 	anims, err := s.GetAnimations(sampleAnimMetas)
// 	assert.NoError(t, err)
// 	t.Fatal(anims)
// }

// func TestSpriteMetaManager(t *testing.T) {
// 	s, err := render.NewSpriteMetaManager(csv.NewReader(strings.NewReader(sampleSpriteManifest)))
// 	assert.NoError(t, err)
// 	t.Fatal(s)
// }

func TestNewSpriteFactoryFromManifests(t *testing.T) {
	factory, err := render.NewSpriteFactoryFromManifests(
		csv.NewReader(strings.NewReader(sampleSheetManifest)),
		csv.NewReader(strings.NewReader(sampleSpriteManifest)),
	)
	assert.NoError(t, err)
	sprite, err := factory.GetSprite("slime1")
	assert.NoError(t, err)
	assert.NotNil(t, sprite)
	// TODO: I spot-checked the sprite with `t.Fatal(sprite)`, but should add real tests
}

// func TestGetSprite(t *testing.T) {
// 	var s render.SpriteGetter

// 	sprite, err := s.GetSprite(id)
// 	assert.NoError(t, err)
// 	assert.Equal(t, *sprite, expected)
// }

const (
	sampleSheetManifest = `name,path,tileX,tileY,sheetX,sheetY
slimes1,./testdata/test_sprite.png,51,54,204,54`
	sampleSpriteManifest = `name,sheet,mode,start,nFrames,dimX,dimY,delay
slime1,slimes1,idle,0,4,51,54,2`
	sampleSpriteID = "slime1"
	sampleSheetID  = "slimes1"
)

var (
	sampleAnimMetas []render.AnimMeta = []render.AnimMeta{
		{
			Mode:       render.NoAnimation,
			Source:     sampleSheetID,
			Start:      0,
			NFrames:    4,
			FrameDelay: 2,
		},
	}
)
