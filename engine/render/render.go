package render

import (
	"encoding/csv"
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: Make a sprite sheet powered implementation of AnimationGetter

type (
	// Animation : a sequence of frames for a sprite to draw. Resist direct member updates; they're
	// exported for serializer.
	Animation struct {
		*Tile
		// Frames : The sequence of images in this animation
		Frames []*ebiten.Image
		// CurrentFrame : which frame in the sequence are we on?
		CurrentFrame int
		// FrameDelay : how much longer should we stay on this frame?
		FrameDelay int
		// FrameDelayMax : how long should we stay on any given frame?
		FrameDelayMax int
	}
	// AnimationMode : An identifier for an animation registered to an Animatable.
	AnimationMode string
	// AnimMeta : Animation metadata for retrieval
	AnimMeta struct {
		Mode       AnimationMode
		Source     SourceImageID
		Start      int
		NFrames    int
		FrameDelay int
	}
	// BasicSprite : A collection of ready-to-render animations
	BasicSprite struct {
		*Animation
		// registeredAnimations : animations available to this sprite
		registeredAnimations map[AnimationMode]Animation
	}
	// Point : A rank-2 vector
	Point struct {
		X float64
		Y float64
	}
	// SpriteID : An identifier for a registered sprite
	SpriteID string
	// SpriteFactory : A simple pipe from SpriteMetaGetter to AnimationGetter to assmeble a Sprite
	SpriteFactory struct {
		sourceImageGetter SpriteSheetGetter
		spriteMetaGetter  SpriteMetaGetter
	}
	// Animator : Triggers a registered animation. Returns the number of frames in a loop
	Animator interface {
		Animate(AnimationMode) (int, error)
	}
	// Drawer : Able to hook into the game engine's animation loop
	Drawer interface {
		Draw(*ebiten.Image)
		Update() error
	}
	// Sprite : An animatible, transformable entity
	Sprite interface {
		Animator
		Transformer
		Drawer
	}
	// SpriteGetter : Give me a Sprite!
	SpriteGetter interface {
		// GetSprite : Get a drawable, animatable, transformable Sprite
		GetSprite(spriteID SpriteID) (Sprite, error)
	}
	// SpriteMetaGetter : Provides metadata for Sprite retrieval
	SpriteMetaGetter interface {
		// GetSpriteMeta : Get Sprite metadata from an ID
		GetSpriteMeta(SpriteID) (SpriteMeta, error)
	}
	// Tile : Static drawable
	Tile struct {
		// image : the thing to render
		image *ebiten.Image
		// position : the top-left of where this should render
		position Point
		// dims : how big should this render?
		dims image.Point
	}
	// Transformer : A handle for modifying scale, location, and rotation
	Transformer interface {
		// Scale : Rescale the transform to the given dimensions
		Scale(Point) error
		// Translate : Move the transform to a new position
		Translate(Point) error
		// Skew : Like rotation, but not!
		// Skew() error // TODO:
		// Rotate : Rotate the transform about the center
		// Rotate() error // TODO: Units? Degrees or radians?
		// GetDims : Get the current post-scale dimensions
		GetDims() image.Point
		// GetPosition : Get the current top-left pixel of the transform
		GetPosition() Point
	}
)

// NoAnimation : A placeholder AnimationMode for when no animation is occurring
const NoAnimation AnimationMode = "no_animation"

// NewSpriteFactory : Create a new pipeline from SpriteID to Sprite
func NewSpriteFactory(a SpriteSheetGetter, s SpriteMetaGetter) (*SpriteFactory, error) {
	return &SpriteFactory{
		sourceImageGetter: a,
		spriteMetaGetter:  s,
	}, nil
}

// NewSpriteFactoryFromManifests : Create a new Sprite generator with simple args
func NewSpriteFactoryFromManifests(sheetManifest, spriteManifest *csv.Reader) (*SpriteFactory, error) {
	sheetManager, err := NewSpriteSheetFiles(sheetManifest)
	if err != nil {
		return nil, err
	}
	spriteMetaManager, err := NewSpriteMetaManager(spriteManifest)
	if err != nil {
		return nil, err
	}
	return NewSpriteFactory(sheetManager, spriteMetaManager)
}

// NewTile : Make a new fixed image renderable
func NewTile(i *ebiten.Image) Tile {
	return Tile{image: i}
}

// Update : Hook for the engine's tick function
func (a *Animation) Update() error {
	if len(a.Frames) == 1 {
		a.image = a.Frames[0]
		return nil
	}
	a.FrameDelay--
	if a.FrameDelay < 0 {
		a.FrameDelay = a.FrameDelayMax
		a.CurrentFrame++
		if a.CurrentFrame >= len(a.Frames) {
			a.CurrentFrame = 0
		}
	}
	a.image = a.Frames[a.CurrentFrame]
	return nil
}

// Animate : Switch to a new Animation
func (b *BasicSprite) Animate(mode AnimationMode) (int, error) {
	anim, ok := b.registeredAnimations[mode]
	if !ok {
		return 0, fmt.Errorf("animation %v could not be found for sprite %v", mode, b)
	}
	anim.CurrentFrame = 0
	anim.FrameDelay = anim.FrameDelayMax
	b.Animation = &anim
	return len(anim.Frames) * anim.FrameDelay, nil
}

// Dist : The distance to the other Point.
func (p *Point) Dist(p2 Point) float64 {
	return math.Sqrt(math.Pow(p.Y-p2.Y, 2) + math.Pow(p.X-p2.X, 2))
}

// AddVec : And a vector of magnitude mag, pointing towards dir
func (p *Point) AddVec(mag float64, dir Point) Point {
	scalar := mag / math.Sqrt(dir.X*dir.X+dir.Y*dir.Y)
	return Point{(dir.X - p.X) * scalar, (dir.Y - p.Y) * scalar}
}

// Lerp : Linear interpolation. Returns an Update hook that moves p to the given endpoint over the
// given number of ticks
func (p *Point) Lerp(Point, int) func() error {
	panic("not implemented") // TODO: Implement
	return func() error {
		return nil
	}
}

// Draw : Engine Draw hook
func (t *Tile) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// TODO: figure out rescaling
	op.GeoM.Translate(t.position.X, t.position.Y)
	screen.DrawImage(t.image, op)
}

// Update : Engine Update hook
func (t *Tile) Update() error { return nil }

// Scale : Rescale the transform to the given dimensions
func (t *Tile) Scale(_ Point) error {
	panic("not implemented") // TODO: Implement
}

// Translate : Move the transform to a new position
func (t *Tile) Translate(_ Point) error {
	panic("not implemented") // TODO: Implement
}

// Skew : Like rotation, but not!
// TODO:
func (t *Tile) Skew() error {
	panic("not implemented") // TODO: Implement
}

// Rotate : Rotate the transform about the center
// TODO: Units? Degrees or radians?
func (t *Tile) Rotate() error {
	panic("not implemented") // TODO: Implement
}

// GetDims : Get the current post-scale dimensions
func (t *Tile) GetDims() image.Point {
	return t.dims
}

// GetPosition : Get the current top-left pixel of the transform
func (t *Tile) GetPosition() Point {
	return t.position
}

// GetSprite : Get a Sprite using the provided SpriteMetaGetter and AnimationGetter
func (s *SpriteFactory) GetSprite(id SpriteID) (Sprite, error) {
	meta, err := s.spriteMetaGetter.GetSpriteMeta(id)
	if err != nil {
		return nil, err
	}
	anims, err := s.GetAnimations(s.sourceImageGetter, meta.Anims)
	if err != nil {
		return nil, err
	}
	result := &BasicSprite{
		Animation: &Animation{
			Tile: &Tile{
				position: Point{0, 0},
			},
		},
		registeredAnimations: anims,
	}
	return result, nil
}

// GetAnimations : Retrieve animations from sprite sheets
func (s *SpriteFactory) GetAnimations(source SpriteSheetGetter, metas []AnimMeta) (map[AnimationMode]Animation, error) {
	batches := map[SourceImageID][]AnimMeta{}
	for _, meta := range metas {
		if batch, ok := batches[meta.Source]; ok {
			batch = append(batch, meta)
			continue
		}
		batches[meta.Source] = []AnimMeta{meta}
	}
	result := map[AnimationMode]Animation{}
	for sheetID, metas := range batches {
		// RFE: Is it worthwhile to use sync map writing to parallelize this?
		sheet, err := source.GetSpriteSheet(sheetID)
		if err != nil {
			return nil, err
		}
		for _, meta := range metas {
			anim, err := sheet.ExtractAnimation(meta)
			if err != nil {
				return nil, err
			}
			result[meta.Mode] = anim
		}
	}
	return result, nil
}
