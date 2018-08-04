package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	*Entity
}

func NewPlayer(x, y int32) *Player {
	stateMap := make(map[int]EntityState)
	stateMap[ENTITY_STATE_IDLE] = EntityState{ Asset: "Player/Idle", FlipHorizontal: false, FlipVertical: false, InheritFlip: true, FrameRate: 200 }
	stateMap[ENTITY_STATE_LEFT] = EntityState{ Asset: "Player/Run", FlipHorizontal: true, FlipVertical: false, InheritFlip: false, FrameRate: 60 }
	stateMap[ENTITY_STATE_RIGHT] = EntityState{ Asset: "Player/Run", FlipHorizontal: false, FlipVertical: false, InheritFlip: false, FrameRate: 60 }
	stateMap[ENTITY_STATE_JUMP] = EntityState{ Asset: "Player/Fall-Jump-WallJ/Jump", FlipHorizontal: false, FlipVertical: false, InheritFlip: true, FrameRate: 200 }
	return &Player {NewEntity(stateMap, x, y) }
}

func (p *Player) Update() {
	p.HandleInput()
	p.Entity.UpdateState()
	p.Entity.Update()
}

func (p *Player) HandleInput() {
	p.Entity.MoveRight = inputManager.KeysHeld[sdl.K_RIGHT]
	p.Entity.MoveLeft = inputManager.KeysHeld[sdl.K_LEFT]
	if inputManager.KeysHeld[sdl.K_SPACE] {
		p.Entity.Jump()
	}
}
