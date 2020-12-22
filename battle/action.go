package battle

import "github.com/hajimehoshi/ebiten/v2"

type (
	// actioner : Performs a combat action
	actioner interface {
		// action : Perform a combat action on some number of targets
		action(b *BattleScene) error
	}
	queuedAction struct {
		act func(*BattleScene, *pawn, ...*pawn) error
		icon ebiten.Image
		executor *pawn
		targets []*pawn
	}
)

func (q *queuedAction) action(b *BattleScene) error {
	return q.act(b, q.executor, q.targets...)
}
