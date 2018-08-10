package main

type Heart struct {
	*Entity
}

func NewHeart(x, y int32) *Heart {
	stateMap := make(map[int]EntityState)
	stateMap[ENTITY_STATE_IDLE] = EntityState{Asset: "Items/Heart/Pick heart", FlipVertical: false, FlipHorizontal: false, InheritFlip: true, FrameRate: 150 }
	ent := NewEntity(stateMap, x, y)
	return &Heart{ent }
}
