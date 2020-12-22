package render_test

import (
	"spiders/engine"
	"testing"

	"github.com/stretchr/testify/assert"
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
