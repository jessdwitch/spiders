package render

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: Make a sprite sheet powered implementation of AnimationGetter

type (
	// Animation : a sequence of frames for a sprite to draw. Resist direct member updates; they're
	// exported for serializer.
	Animation struct {
		Tile
		// Frames : The sequence of images in this animation
		Frames []*ebiten.Image
		// CurrentFrame : which frame in the sequence are we on?
		CurrentFrame int
		// FrameDelay : how much longer should we stay on this frame?
		FrameDelay int
		// FrameDelayMax : how long should we stay on any given frame?
		FrameDelayMax int
		// Dims : how big is this sprite?
		Dims Point
	}
	// AnimationMode : An identifier for an animation registered to an Animatable.
	AnimationMode string
	// AnimMeta : Animation metadata for retrieval
	AnimMeta struct {
		mode    AnimationMode
		source  SourceImageID
		start   int
		nFrames int
	}
	// BasicSprite : A collection of ready-to-render animations
	BasicSprite struct {
		Animation
		// registeredAnimations : animations available to this sprite
		registeredAnimations map[AnimationMode]Animation
	}
	// Point : A rank-2 vector
	Point struct {
		X float64
		Y float64
	}
	// SourceImageID : An identifier for a registered source image provider
	SourceImageID int
	// SpriteID : An identifier for a registered sprite
	SpriteID int
	// SpriteMeta :
	SpriteMeta struct {
		InitialDims Point
	}
	// animationCycle : How should this animation cycle?
	animationCycle int
	// SpriteManager : A simple pipe from SpriteMetaGetter to AnimationGetter to assmeble a Sprite
	SpriteManager struct {
		animationGetter  AnimationGetter
		spriteMetaGetter SpriteMetaGetter
		// spriteAnimations map[SpriteID][]AnimMeta
		// spriteSheets     map[SourceImageID]sheetMeta
	}
	sheetMeta struct {
		fp       string
		sheetDim image.Point
		tileDim  image.Point
		nXTiles  int
		nYTiles  int
	}
	// AnimationGetter : Assemble an animation from it's metadata
	AnimationGetter interface {
		GetAnimations([]AnimMeta) (map[AnimationMode]Animation, error)
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
		image    *ebiten.Image
		position Point
	}
	// Transformer : A handle for modifying scale, location, and rotation
	Transformer interface {
		// Scale : Rescale the transform to the given dimensions
		Scale(Point) error
		// Translate : Move the transform to a new position
		Translate(Point) error
		// Skew : Like rotation, but not!
		Skew() error // TODO:
		// Rotate : Rotate the transform about the center
		Rotate() error // TODO: Units? Degrees or radians?
		// GetDims : Get the current post-scale dimensions
		GetDims() Point
		// GetPosition : Get the current top-left pixel of the transform
		GetPosition() Point
	}
)

// NoAnimation : A placeholder AnimationMode for when no animation is occurring
const NoAnimation AnimationMode = "no_animation"

// NewSpriteManager : Create a new pipeline from SpriteID to Sprite
func NewSpriteManager(a AnimationGetter, s SpriteMetaGetter) (*SpriteManager, error) {
	return &SpriteManager{
		animationGetter:  a,
		spriteMetaGetter: s,
		// spriteAnimations: make(map[SpriteID][]AnimMeta),
		// spriteSheets:     make(map[SourceImageID]sheetMeta),
	}, nil
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
	b.Animation = anim
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
func (t *Tile) GetDims() Point {
	panic("not implemented") // TODO: Implement
}

// GetPosition : Get the current top-left pixel of the transform
func (t *Tile) GetPosition() Point {
	panic("not implemented") // TODO: Implement
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

func (s *SpriteManager) GetSpriteAnimations(metas []AnimMeta) (map[AnimationMode]Animation, error) {
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
		sheet := s.spriteSheets[sheetID]
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

// GetSprite : Get a Sprite using the provided SpriteMetaGetter and AnimationGetter
func (s *SpriteManager) GetSprite(id SpriteID) (*BasicSprite, error) {
	meta, err := s.spriteMetaGetter.GetSpriteMeta(id)
	if err != nil {
		return nil, err
	}
	anims, err := s.animationGetter.GetAnimations(meta.AnimMetas)
	if err != nil {
		return nil, err
	}
	result := &BasicSprite{
		Animation: Animation{
			Tile: Tile{
				position: Point{0, 0},
			},
		},
		registeredAnimations: anims,
	}
	return result, nil
}

func (s *sheetMeta) iToRect(i int) (image.Rectangle, error) {
	if i < 0 {
		return image.Rectangle{}, fmt.Errorf("requested index %d is negative for sheet %s", i, s.fp)
	}
	if i >= s.nXTiles*s.nYTiles {
		return image.Rectangle{}, fmt.Errorf("tile index %d must be less than %d for sheet %s", i, s.nXTiles*s.nYTiles, s.fp)
	}
	x := (i % s.nXTiles) * s.tileDim.X
	y := (i / s.nXTiles) * s.tileDim.Y
	return image.Rect(x, y, x+s.tileDim.X, y+s.tileDim.Y), nil
}
