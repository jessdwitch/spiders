package render

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type (
	// AnimationGetter : Assemble an animation from it's metadata
	AnimationGetter interface {
		GetAnimations([]AnimMeta) (map[AnimationMode]Animation, error)
	}
	// SourceImageID : An identifier for a registered source image provider
	SourceImageID string
	// SpriteSheetManager : A provider of sprite sheets
	SpriteSheetManager map[SourceImageID]sheetMeta
	sheetMeta          struct {
		fp       string
		sheetDim image.Point
		tileDim  image.Point
		nTiles   image.Point
	}
)

// NewSpriteSheetManager : Get a new AnimationGetter from a sprite sheet manifest
func NewSpriteSheetManager(manifest csv.Reader) (SpriteSheetManager, error) {
	result := SpriteSheetManager(make(map[SourceImageID]sheetMeta))
	// strip header
	_, err := manifest.Read()
	if err == io.EOF {
		return nil, fmt.Errorf("sheet manifest is empty")
	}
	if err != nil {
		return nil, err
	}
	for {
		record, err := manifest.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		err = result.processManifestCsvRecord(record)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// GetAnimations : Retrieve animations from sprite sheets
func (s SpriteSheetManager) GetAnimations(metas []AnimMeta) (map[AnimationMode]Animation, error) {
	batches := map[SourceImageID][]AnimMeta{}
	for _, meta := range metas {
		if batch, ok := batches[meta.source]; ok {
			batch = append(batch, meta)
		} else {
			batches[meta.source] = []AnimMeta{meta}
		}
	}
	result := map[AnimationMode]Animation{}
	for sheetID, metas := range batches {
		// RFE: Is it worthwhile to use sync map writing to parallelize this?
		sheet := s[sheetID]
		for _, meta := range metas {
			anim, err := sheet.getAnimation(meta)
			if err != nil {
				return nil, err
			}
			result[meta.mode] = anim
		}
	}
	return result, nil
}

func (s SpriteSheetManager) processManifestCsvRecord(record []string) error {
	// record: name, path, tileX, tileY
	var err error
	meta := sheetMeta{
		fp:      record[1],
		tileDim: image.Point{},
	}
	meta.tileDim.X, err = strconv.Atoi(record[3])
	if err != nil {
		return err
	}
	meta.tileDim.Y, err = strconv.Atoi(record[4])
	if err != nil {
		return err
	}

	meta.sheetDim.X, err = strconv.Atoi(record[5])
	if err != nil || meta.sheetDim.X == 0 {
		f, err := os.Open(meta.fp)
		if err != nil {
			return err
		}
		defer f.Close()

		i, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		meta.sheetDim.X = i.Bounds().Max.X
		meta.sheetDim.Y = i.Bounds().Max.Y
	} else {
		meta.sheetDim.Y, err = strconv.Atoi(record[6])
		if err != nil {
			return err
		}
	}

	meta.nTiles.X = meta.sheetDim.X / meta.tileDim.X
	meta.nTiles.Y = meta.sheetDim.Y / meta.tileDim.Y

	s[SourceImageID(record[0])] = meta

	return nil
}

func (s *sheetMeta) getAnimation(meta AnimMeta) (Animation, error) {
	sheetImg, _, err := ebitenutil.NewImageFromFile(s.fp)
	if err != nil {
		return Animation{}, err
	}
	result := Animation{
		Frames: []*ebiten.Image{},
	}
	for i := 0; i < meta.nFrames; i++ {
		r, err := s.iToRect(meta.start + i)
		if err != nil {
			return Animation{}, err
		}
		frame := sheetImg.SubImage(r).(*ebiten.Image)
		result.Frames = append(result.Frames, frame)
	}
	return result, nil
}

func (s *sheetMeta) iToRect(i int) (image.Rectangle, error) {
	if i < 0 {
		return image.Rectangle{}, fmt.Errorf("requested index %d is negative for sheet %s", i, s.fp)
	}
	if i >= s.nTiles.X*s.nTiles.Y {
		return image.Rectangle{},
			fmt.Errorf("tile index %d must be less than %d for sheet %s", i, s.nTiles.X*s.nTiles.Y, s.fp)
	}
	x := (i % s.nTiles.X) * s.tileDim.X
	y := (i / s.nTiles.Y) * s.tileDim.Y
	return image.Rect(x, y, x+s.tileDim.X, y+s.tileDim.Y), nil
}
