package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	*Entity
	CanJump bool
	Subscriptions []Subscription
}

func NewPlayer(x, y int32) *Player {
	stateMap := make(map[int]EntityState)
	stateMap[ENTITY_STATE_IDLE] = EntityState{ Asset: "Player/Idle", FlipHorizontal: false, FlipVertical: false, InheritFlip: true, FrameRate: 200 }
	stateMap[ENTITY_STATE_LEFT] = EntityState{ Asset: "Player/Run", FlipHorizontal: true, FlipVertical: false, InheritFlip: false, FrameRate: 60 }
	stateMap[ENTITY_STATE_RIGHT] = EntityState{ Asset: "Player/Run", FlipHorizontal: false, FlipVertical: false, InheritFlip: false, FrameRate: 60 }
	stateMap[ENTITY_STATE_JUMP] = EntityState{ Asset: "Player/Fall-Jump-WallJ/Jump", FlipHorizontal: false, FlipVertical: false, InheritFlip: true, FrameRate: 200 }

	p := &Player { Entity: NewEntity(stateMap, x, y), CanJump: true }

	jumpSub := inputManager.RegisterKeyListener(sdl.K_SPACE, "pressed", func() {
		if !p.Entity.IsJumping && p.CanJump{
			p.Entity.Jump()
			p.CanJump = false
		}
	})
	jumpEndSub := inputManager.RegisterKeyListener(sdl.K_SPACE, "released", func() {
		p.CanJump = true
	})

	p.Subscriptions = append(p.Subscriptions, jumpSub, jumpEndSub)
	return p
}

func (p *Player) Update() {
	p.HandleInput()
	p.Entity.UpdateState()
	p.Entity.Update()
}

func (p *Player) HandleInput() {
	p.Entity.MoveRight = inputManager.KeysHeld[sdl.K_RIGHT] || inputManager.KeysHeld[sdl.K_d]
	p.Entity.MoveLeft = inputManager.KeysHeld[sdl.K_LEFT] || inputManager.KeysHeld[sdl.K_a]
}

func (p *Player) Cleanup() {
	for _, sub := range p.Subscriptions {
		sub.Unsubscribe()
	}
}
