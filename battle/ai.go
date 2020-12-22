package battle

type (
	// routiner : Given the current battle state, choose a course an action
	routiner interface {
		routine(b *BattleScene, reason routineReason) (actioner, error)
	}
	routineReason int
	enemy struct {
		id int
		name string
	}
)

func (e *enemy) routine(b *BattleScene, event turnEvent) (actioner, error) {
	return nil, nil
}

func newEnemiesFromIDs(ids []int) ([]pawn, error) {
	var err error
	result := make([]pawn, len(ids))
	for i, id := range ids {
		// TODO: Depending on how we store enemy data, this may want batching (as in reducing calls to Open by caching them in a map)
		result[i], err = newEnemyFromID(id)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func newEnemyFromID(id int) (pawn, error) {
	return pawn{}, nil
}
