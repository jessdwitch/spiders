package engine

import (
	"github.com/jessdwitch/spiders/engine/render"

	"fmt"
)

type (
	PlayerParty struct {
		ActiveMembers []Character
		// AllMembers []Character
	}
	Character struct {
		Name          string
		MaxHealth     int
		CurrentHealth int
		sprites       map[string]render.SpriteID
	}
)

func (c *Character) GetSprite(s render.SpriteGetter, context string) (render.Sprite, error) {
	id, ok := c.sprites[context]
	if !ok {
		return nil, fmt.Errorf("character %s does not have sprite for context %s", c.Name, context)
	}
	return s.GetSprite(id)
}
