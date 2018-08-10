package main

type Coin struct {
	*Entity
}

func NewCoin(x, y int32) *Coin {
	stateMap := make(map[int]EntityState)
	stateMap[ENTITY_STATE_IDLE] = EntityState{Asset: "Items/Coin/Spin", FlipVertical: false, FlipHorizontal: false, InheritFlip: true, FrameRate: 150 }
	ent := NewEntity(stateMap, x, y)
	return &Coin{ent }
}
