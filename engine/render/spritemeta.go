package render

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type (
	// SpriteMeta :
	SpriteMeta struct {
		InitialDims Point
		Anims       []AnimMeta
	}
	// SpriteMetaManager :
	SpriteMetaManager map[SpriteID]SpriteMeta
)

// NewSpriteMetaManager :
func NewSpriteMetaManager(manifest *csv.Reader) (SpriteMetaManager, error) {
	result := SpriteMetaManager(make(map[SpriteID]SpriteMeta))
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

func (s SpriteMetaManager) processManifestCsvRecord(record []string) error {
	// record: name, sheet, mode, start, nFrames, dimX, dimY, delay
	var err error
	meta, ok := s[SpriteID(record[0])]
	if !ok {
		// TODO: Rescale for screen size
		dims := Point{}
		dims.X, err = strconv.ParseFloat(record[5], 64)
		dims.Y, err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			return err
		}
		meta = SpriteMeta{
			dims,
			[]AnimMeta{},
		}
	}
	start, err := strconv.Atoi(record[3])
	if err != nil {
		return err
	}
	nFrames, err := strconv.Atoi(record[4])
	if err != nil {
		return err
	}
	delay, err := strconv.Atoi(record[7])
	if err != nil {
		return err
	}
	anim := AnimMeta{
		Mode:       AnimationMode(record[2]),
		Source:     SourceImageID(record[1]),
		Start:      start,
		NFrames:    nFrames,
		FrameDelay: delay,
	}
	meta.Anims = append(meta.Anims, anim)
	s[SpriteID(record[0])] = meta
	return nil
}

// GetSpriteMeta : Get Sprite metadata from an ID
func (s SpriteMetaManager) GetSpriteMeta(id SpriteID) (SpriteMeta, error) {
	if meta, ok := s[id]; ok {
		return meta, nil
	}
	return SpriteMeta{}, fmt.Errorf("id %s not found in sprite metas", id)
}
