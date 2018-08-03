package main

type Player struct {
	Entity
}

func NewPlayer() *Entity {
	return NewEntity("Player/Idle")
}

func (p *Player) Update() {
	p.Entity.Update()
}

func (p *Player) SetState(state int) {
	var newState string
	// state := p.Entity.CurrentState
	if state == ENTITY_STATE_RIGHT {
		newState = "Player/Run"
		p.Entity.Image.FlipHorizontal = false
	}
	if state == ENTITY_STATE_LEFT {
		newState = "Player/Run"
		p.Entity.Image.FlipHorizontal = true
	}
	if state == ENTITY_STATE_IDLE {
		newState = "Player/Idle"
	}
	if state == ENTITY_STATE_JUMP {
		newState = "Player/Fall-Jump-WallJ/Jump"
	}
	p.Entity.Image.SetAsset(newState)
}