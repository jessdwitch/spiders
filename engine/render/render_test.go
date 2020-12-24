package render_test

import (
	"testing"

	"github.com/jessdwitch/spiders/engine"

	"github.com/stretchr/testify/assert"

	_ "image/png"
)

func makeSpriteManager() (*engine.SpriteGetter, engine.SpriteID, engine.BasicSprite) {
	s := &engine.NewSpriteGetter()

}

func TestGetSprite(t *testing.T) {
	var s engine.SpriteGetter
	s, id, expected := makeSpriteManager()
	sprite, err := s.GetSprite(id)
	assert.NoError(t, err)
	assert.Equal(t, *sprite, expected)
}
