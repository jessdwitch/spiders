package render

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type (
	// AnimationGetter : Assemble an animation from it's metadata
	AnimationGetter interface {
		GetAnimations([]AnimMeta) (map[AnimationMode]Animation, error)
	}
	// SpriteSheetGetter : Get a sprite source image
	SpriteSheetGetter interface {
		GetSpriteSheet(SourceImageID) (*SpriteSheet, error)
	}
	// SourceImageID : An identifier for a registered source image provider
	SourceImageID string
	// SpriteSheet : An image with sprite extraction details
	SpriteSheet struct {
		SourceImage *ebiten.Image
		SheetID     SourceImageID
		SheetDim    image.Point
		TileDim     image.Point
		NTiles      image.Point
	}
	// SpriteSheetFiles : A provider of sprite sheets from file sources
	SpriteSheetFiles map[SourceImageID]sheetFileMeta
	sheetFileMeta    struct {
		fp       string
		sheetDim image.Point
		tileDim  image.Point
		nTiles   image.Point
	}
)

// NewSpriteSheetFiles : Get a new AnimationGetter from a sprite sheet manifest
func NewSpriteSheetFiles(manifest *csv.Reader) (SpriteSheetFiles, error) {
	result := SpriteSheetFiles(make(map[SourceImageID]sheetFileMeta))
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

func (s SpriteSheetFiles) processManifestCsvRecord(record []string) error {
	// record: name, path, tileX, tileY, sheetX, sheetY
	var err error
	meta := sheetFileMeta{
		fp:      record[1],
		tileDim: image.Point{},
	}
	meta.tileDim.X, err = strconv.Atoi(record[2])
	if err != nil {
		return err
	}
	meta.tileDim.Y, err = strconv.Atoi(record[3])
	if err != nil {
		return err
	}

	meta.sheetDim.X, err = strconv.Atoi(record[4])
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
		meta.sheetDim.Y, err = strconv.Atoi(record[5])
		if err != nil {
			return err
		}
	}

	meta.nTiles.X = meta.sheetDim.X / meta.tileDim.X
	meta.nTiles.Y = meta.sheetDim.Y / meta.tileDim.Y

	s[SourceImageID(record[0])] = meta

	return nil
}

// GetSpriteSheet : Get a registered SpriteSheet
func (s SpriteSheetFiles) GetSpriteSheet(id SourceImageID) (*SpriteSheet, error) {
	if meta, ok := s[id]; ok {
		return meta.GetSpriteSheet(id)
	}
	return nil, fmt.Errorf("sheet %s not found", id)
}

// ExtractAnimation : Extract a series of frames from this Sprite Sheet
func (s *SpriteSheet) ExtractAnimation(meta AnimMeta) (Animation, error) {
	result := Animation{
		Tile:          &Tile{dims: s.TileDim},
		Frames:        []*ebiten.Image{},
		FrameDelayMax: meta.FrameDelay,
	}
	for i := 0; i < meta.NFrames; i++ {
		r, err := s.iToRect(meta.Start + i)
		if err != nil {
			return Animation{}, err
		}
		frame := s.SourceImage.SubImage(r).(*ebiten.Image)
		result.Frames = append(result.Frames, frame)
	}
	return result, nil
}

func (s *SpriteSheet) iToRect(i int) (image.Rectangle, error) {
	if i < 0 {
		return image.Rectangle{}, fmt.Errorf("requested index %d is negative for sheet %s", i, s.SheetID)
	}
	if i >= s.NTiles.X*s.NTiles.Y {
		return image.Rectangle{},
			fmt.Errorf("tile index %d must be less than %d for sheet %s", i, s.NTiles.X*s.NTiles.Y, s.SheetID)
	}
	x := (i % s.NTiles.X) * s.TileDim.X
	y := (i / s.NTiles.Y) * s.TileDim.Y
	return image.Rect(x, y, x+s.TileDim.X, y+s.TileDim.Y), nil
}

func (s *sheetFileMeta) GetSpriteSheet(id SourceImageID) (*SpriteSheet, error) {
	f, err := os.Open(s.fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return &SpriteSheet{
		SourceImage: ebiten.NewImageFromImage(img),
		SheetID:     id,
		SheetDim:    s.sheetDim,
		TileDim:     s.tileDim,
		NTiles:      s.nTiles,
	}, nil
}
