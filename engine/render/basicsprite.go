package render

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type (
	// BasicSprite : Implements Sprite
	BasicSprite struct {
		// position: Where to render the top-left of the sprite
		position Point
		// dims : The width and height of the Sprite
		dims Point
		// currentMode : What animation is currently running?
		currentMode AnimationMode
		// registeredAnimations : animations available to this sprite
		registeredAnimations map[AnimationMode]Animation
	}
)

// Animator implementation

// Animate : Switch to a new Animation
func (b *BasicSprite) Animate(mode AnimationMode) (int, error) {
	b.currentMode = mode
	anim, ok := b.registeredAnimations[mode]
	if !ok {
		return 0, fmt.Errorf("animation %v could not be found for sprite %v", mode, b)
	}
	anim.CurrentFrame = 0
	anim.FrameDelay = anim.FrameDelayMax
	return len(anim.Frames) * anim.FrameDelay, nil
}

// Transform implementation

// Scale : Rescale the transform to the given dimensions
func (b *BasicSprite) Scale(_ Point) error {
	panic("not implemented") // TODO: Implement
}

// Translate : Move the transform to a new position
func (b *BasicSprite) Translate(_ Point) error {
	panic("not implemented") // TODO: Implement
}

// Lerp : Linear interpolation. Move to the given endpoint over the given number of ticks
func (b *BasicSprite) Lerp(_ Point, _ int) error {
	panic("not implemented") // TODO: Implement
}

// Skew : Like rotation, but not!
// TODO:
func (b *BasicSprite) Skew() error {
	panic("not implemented") // TODO: Implement
}

// Rotate : Rotate the transform about the center
// TODO: Units? Degrees or radians?
func (b *BasicSprite) Rotate() error {
	panic("not implemented") // TODO: Implement
}

// GetDims : Get the current post-scale dimensions
func (b *BasicSprite) GetDims() Point {
	panic("not implemented") // TODO: Implement
}

// GetPosition : Get the current top-left pixel of the transform
func (b *BasicSprite) GetPosition() Point {
	panic("not implemented") // TODO: Implement
}

// Sprite interface additional implementations

// Update : Hook for the engine's tick function
func (b *BasicSprite) Update() error {
	anim := b.registeredAnimations[b.currentMode]
	if len(anim.Frames) == 1 {
		return nil
	}
	anim.FrameDelay--
	if anim.FrameDelay < 0 {
		anim.FrameDelay = anim.FrameDelayMax
		anim.CurrentFrame++
		if anim.CurrentFrame >= len(anim.Frames) {
			anim.CurrentFrame = 0
		}
	}
	return nil
}

// Draw : Draw this Sprite onto the given image
func (b *BasicSprite) Draw(screen *ebiten.Image) {
	anim := b.registeredAnimations[b.currentMode]
	op := &ebiten.DrawImageOptions{}
	// TODO: figure out rescaling
	op.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(anim.Frames[anim.CurrentFrame], op)
}
