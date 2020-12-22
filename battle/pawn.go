package battle

import (
	"github.com/jessdwitch/spiders/engine"
	"github.com/jessdwitch/spiders/engine/render"

	"github.com/hajimehoshi/ebiten/v2"
)

type (
	pawn struct {
		name          string
		maxHealth     int
		currentHealth int
		statuses      []status
		sprite        render.Sprite
		action        queuedAction
	}
	pawns          []pawn
	statusResolver interface {
		// resolveStatus : Provide an action to carry out, and a status to replace with after resolution
		resolveStatus(turnEvent) (actioner, status, error)
	}
	status interface {
		statusResolver
		Draw(*ebiten.Image)
	}
)

func newPawnsFromParty(party engine.PlayerParty) (pawns, error) {
	result := make(pawns, len(party.ActiveMembers))
	for i, p := range party.ActiveMembers {
		result[i] = pawn{
			name:          p.Name,
			maxHealth:     p.MaxHealth,
			currentHealth: p.CurrentHealth,
			statuses:      []status{},
		}
	}
	return result, nil
}

func (p *pawn) resolveStatuses(b *BattleScene) error {
	newStatuses := []status{}
	for _, s := range p.statuses {
		act, sNew, err := s.resolveStatus(b.state)
		if err != nil {
			return err
		}
		if act != nil {
			if err = act.action(b); err != nil {
				return err
			}
		}
		if sNew != nil {
			newStatuses = append(newStatuses, sNew)
		}
	}
	p.statuses = newStatuses
	return nil
}

func (p *pawns) resolveStatuses(b *BattleScene) error {
	var err error
	for _, pawn := range *p {
		if err = pawn.resolveStatuses(b); err != nil {
			return err
		}
	}
	return nil
}

func (p *pawns) update() error {
	var err error
	for _, pawn := range *p {
		if err = pawn.sprite.Update(); err != nil {
			return err
		}
	}
	return nil
}

// arrange : sets the locations for each of the pawns equidistant on a given line
func (p *pawns) arrange(start, end render.Point) {
	// TODO: Rotate the sprite to match the angle of the line or compute the across of the sprite at that angle
	var totalPawnWidth float64
	for _, pawn := range *p {
		totalPawnWidth += pawn.sprite.GetDims().X
	}
	spacer := (start.Dist(end) - totalPawnWidth) / float64(len(*p)+1)
	for _, pawn := range *p {
		start = start.AddVec(spacer, end)
		pawn.sprite.Translate(start)
		start = start.AddVec(pawn.sprite.GetDims().X, end)
	}
}

func (p *pawns) draw(screen *ebiten.Image) {
	for _, pawn := range *p {
		pawn.sprite.Draw(screen)
	}
}
