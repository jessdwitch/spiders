package battle

import (
	"errors"
	"fmt"

	"github.com/jessdwitch/spiders/deck"
	"github.com/jessdwitch/spiders/engine"
	"github.com/jessdwitch/spiders/engine/render"

	"github.com/hajimehoshi/ebiten/v2"
)

const maxPlayerPawns = 3

type (
	BattleScene struct {
		background     *ebiten.Image
		playerDeck     deck.Deck
		playerHandSize int
		playerPawns    pawns
		actionQueue    []queuedAction
		aiPawns        pawns
		isPlayerTurn   bool
		turnNumber     int
		state          turnEvent
	}
	turnEvent int
)

const (
	invalidEvent turnEvent = iota
	turnStart
	turnInProgress
	turnResolving
	turnEnd
	takeDamage
)

func (t turnEvent) next() turnEvent {
	switch t {
	case turnStart:
		return turnInProgress
	case turnInProgress:
		return turnResolving
	case turnResolving:
		return turnEnd
	case turnEnd:
		return turnStart
	default:
		return invalidEvent
	}
}

// NewBattleScene : Generate a new combat instance
func NewBattleScene(
	gameState *engine.GameState,
	background *ebiten.Image,
	aiIDs []int,
	playerStarts bool,
	playerCards map[deck.CardID]int,
) (*BattleScene, error) {
	// preflight checks
	if len(aiIDs) == 0 {
		return nil, errors.New("battle must have at least one AI pawn")
	}
	b := &BattleScene{}
	playerPawns, err := newPawnsFromParty(gameState.PlayerParty)
	if err != nil {
		return nil, err
	}
	if len(playerPawns) > len(aiIDs) {
		b.actionQueue = make([]queuedAction, len(playerPawns))
	} else {
		b.actionQueue = make([]queuedAction, len(aiIDs))
	}
	b.turnNumber = 1
	b.playerHandSize = 5
	b.state = turnStart
	b.background = background // TODO: Scale to screen
	b.playerDeck, err = deck.NewDeckFromIDs(playerCards)
	if err != nil {
		return nil, err
	}
	b.aiPawns, err = newEnemiesFromIDs(aiIDs)
	if err != nil {
		return nil, err
	}
	b.isPlayerTurn = playerStarts
	playerAxisStart, playerAxisEnd, aiAxisStart, aiAxisEnd := computeAxes(
		gameState.Config.ScreenWidth, gameState.Config.ScreenHeight)
	b.playerPawns.arrange(playerAxisStart, playerAxisEnd)
	b.aiPawns.arrange(aiAxisStart, aiAxisEnd)
	return b, nil
}

func computeAxes(width, height int) (render.Point, render.Point, render.Point, render.Point) {
	panic("not implemented") // TODO: Implement
}

// Update :
func (b *BattleScene) Update(state *engine.GameState) error {
	return nil
}

// Draw : Render the BattleScene, including player and AI pawns, and player UI
func (b *BattleScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(b.background, nil)
	b.playerPawns.draw(screen)
}

func (b *BattleScene) transitionState() error {
	var activePawns pawns
	b.state = b.state.next()
	if b.state == invalidEvent {
		return fmt.Errorf("BattleScene is in an invalid state: %v", *b)
	}
	if b.state == turnStart {
		b.isPlayerTurn = !b.isPlayerTurn
	}
	if b.isPlayerTurn {
		activePawns = b.playerPawns
	} else {
		activePawns = b.aiPawns
	}
	return activePawns.resolveStatuses(b)
}

// func (b *BattleScene) resolveTurnStart() error {
// 	var err error
// 	if b.isPlayerTurn {
// 		if err = b.playerPawns.resolveStatuses(b, turnStart); err != nil {
// 			return err
// 		}
// 	} else {
// 		if err = b.aiPawns.resolveStatuses(b, turnStart); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (b *BattleScene) resolveTurnEnd() error {
// 	var err error
// 	// Handle player turn specific resolution
// 	if b.isPlayerTurn {
// 		for i, p := range b.playerQueue {
// 			if p != nil {
// 				if err = p.action.action(b); err != nil {
// 					return err
// 				}
// 				b.playerQueue[i] = nil
// 			}
// 		}
// 		if err = b.playerPawns.resolveStatuses(b, turnEnd); err != nil {
// 			return err
// 		}
// 		b.playerDeck.DrawCards(b.playerHandSize - len(b.playerDeck.Hand))
// 	} else { // Handle AI turn specific resolution
// 		for _, p := range b.aiPawns {
// 			if err = p.action.action(b); err != nil {
// 				return err
// 			}
// 		}
// 		if err = b.aiPawns.resolveStatuses(b, turnEnd); err != nil {
// 			return err
// 		}
// 		b.turnNumber++
// 	}
// 	b.isPlayerTurn = !b.isPlayerTurn
// 	return nil
// }
