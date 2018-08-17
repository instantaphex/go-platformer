package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MOVEMENT_MASK = (COMPONENT_POSITION|COMPONENT_VELOCITY)
	RENDER_MASK = (COMPONENT_POSITION|COMPONENT_APPEARANCE)
)

type System interface {
	Update(engine *Engine, world *World)
}

func signatureMatches(mask, signature uint64) bool {
	return mask & signature == signature
}

type StateSystem struct {}

func (ss StateSystem) Update(engine *Engine, world *World) {
	var ap *Appearance
	var stateCmp *State

	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_APPEARANCE|COMPONENT_STATE) {
			ap = &(world.appearance[entity])
			stateCmp = &(world.state[entity])

			state := stateCmp.stateMap[stateCmp.currentState]
			ap.name = state.asset
			if !state.inheritFlip {
				ap.flip = state.flip
			}
			ap.inheritFlip = state.inheritFlip
		}
	}
}

type RenderSystem struct {
	signature int64
}
func (rs RenderSystem) Update (engine *Engine, world *World) {
	var pos *Position
	var a *Appearance

	for entity, mask := range world.mask {
		if signatureMatches(mask, RENDER_MASK) {
			pos = &(world.position[entity])
			a = &(world.appearance[entity])

			// apply offset for different sized sprites as well
			// as accounting for camera position
			x := int32(pos.x - engine.Camera.X()) + a.xOffset
			y := int32(pos.y - engine.Camera.Y()) + a.yOffset
			engine.Graphics.DrawPart(engine.Assets.Texture, x, y, a.frame.X, a.frame.Y, a.frame.W, a.frame.H, a.flip)
		}
	}
}

type AnimationSystem struct {}

func (as AnimationSystem) Update(engine *Engine, world *World) {
	var an *Animation
	var ap *Appearance
	var state *State
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_ANIMATION|COMPONENT_APPEARANCE|COMPONENT_STATE) {
			an = &(world.animation[entity])
			ap = &(world.appearance[entity])
			state = &(world.state[entity])

			if an.oldTime + uint32(an.frameRate) > sdl.GetTicks() {
				return
			}
			an.oldTime = sdl.GetTicks()
			an.currentFrame += an.frameInc
			if an.currentFrame >= an.maxFrames {
				an.currentFrame = 0
			}

			// get animation frame array
			state := state.stateMap[state.currentState]
			frames := engine.Assets.Get(state.asset)
			an.maxFrames = len(frames)
			an.frameRate = state.frameRate
			an.frameInc = 1
			if len(frames) == 0 { return }

			// if the animation we're switching to has less
			// frames than the current animation frame, reset it
			if an.currentFrame > len(frames) - 1 {
				an.currentFrame = 0
			}

			// calculate offset to be used in rendering to account
			// for different sized sprite frames
			ap.frame = frames[an.currentFrame]
			ap.xOffset = ap.frame.SourceW - ap.frame.W
			ap.yOffset = ap.frame.SourceH - ap.frame.H
			ap.w = ap.frame.W + ap.xOffset
			ap.h = ap.frame.H + ap.yOffset
		}
	}
}

type InputSystem struct {}
func (cs InputSystem) Update(engine *Engine, world *World) {
	var c *Controllable
	var s *State
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_CONTROLLABLE) {
			c = &(world.controllable[entity])
			s = &(world.state[entity])
			c.moveRight = engine.Input.KeysHeld[sdl.K_RIGHT] || engine.Input.KeysHeld[sdl.K_d]
			c.moveLeft = engine.Input.KeysHeld[sdl.K_LEFT] || engine.Input.KeysHeld[sdl.K_a]
			if c.moveLeft {
				s.currentState = ENTITY_STATE_LEFT
			} else if c.moveRight {
				s.currentState = ENTITY_STATE_RIGHT
			} else {
				s.currentState = ENTITY_STATE_IDLE
			}
			if engine.Input.KeysHeld[sdl.K_SPACE] {
				s.currentState = ENTITY_STATE_JUMP
				c.jumping = true
			}
		}
	}
}

type MovementSystem struct {
	v *Velocity
	pos *Position
	c *Controllable
	a *Appearance
	engine *Engine
}
func (ms *MovementSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_VELOCITY|COMPONENT_POSITION|COMPONENT_APPEARANCE) {
			ms.v = &(world.velocity[entity])
			ms.pos = &(world.position[entity])
			ms.c = &(world.controllable[entity])
			ms.a = &(world.appearance[entity])
			ms.engine = engine

			if !ms.c.moveLeft && !ms.c.moveRight {
				ms.StopMove()
			}
			if ms.c.moveLeft {
				ms.v.accelX = -0.2
			} else if ms.c.moveRight {
				ms.v.accelX = 0.2
			}

			if ms.c.jumping && ms.c.canJump {
				ms.v.speedY = -ms.v.maxSpeedY
				ms.c.canJump = false
			}

			ms.v.accelY = .55
			/*if e.Flags & ENTITY_FLAG_GRAVITY != 0 {
				e.AccelY = .55
			}*/

			ms.v.speedX += ms.v.accelX * engine.FPS.GetSpeedFactor()
			ms.v.speedY += ms.v.accelY * engine.FPS.GetSpeedFactor()

			if ms.v.speedX > ms.v.maxSpeedX { ms.v.speedX = ms.v.maxSpeedX }
			if ms.v.speedX < -ms.v.maxSpeedX { ms.v.speedX = -ms.v.maxSpeedX }
			if ms.v.speedY > ms.v.maxSpeedY { ms.v.speedY = ms.v.maxSpeedY }
			if ms.v.speedY < -ms.v.maxSpeedY { ms.v.speedY = -ms.v.maxSpeedY }

			ms.Move(ms.v.speedX, ms.v.speedY)
		}
	}
}

func (ms *MovementSystem) Move(moveX, moveY float32) {
	if moveX == 0 && moveY == 0 {
		return
	}

	var newX float32 = 0.0
	var newY float32 = 0.0

	moveX *= ms.engine.FPS.GetSpeedFactor()
	moveY *= ms.engine.FPS.GetSpeedFactor()

	// if x or y are not zero, start collision checking
	// at the current x + speed factor and current y + speed factor
	if moveX != 0 {
		if moveX >= 0 {
			newX = ms.engine.FPS.GetSpeedFactor()
		} else {
			newX = -ms.engine.FPS.GetSpeedFactor()
		}
	}

	if moveY != 0 {
		if moveY >= 0 {
			newY = ms.engine.FPS.GetSpeedFactor()
		} else {
			newY = -ms.engine.FPS.GetSpeedFactor()
		}
	}

	// starting at speed factor (+ x, + y) keep iterating over each new desired
	// move, checking for collisions allong the way
	for {
		// check for collision at desired x position
		if ms.PosValid(int32(ms.pos.x + newX), int32(ms.pos.y)) {
			// no collision, update position
			ms.pos.x += newX
		} else {
			// collision deteced, set speed to zero so we don't allow the movement
			ms.v.speedX = 0
		}

		// check for collision at desired y
		if ms.PosValid(int32(ms.pos.x), int32(ms.pos.y + newY)) {
			// no collision, update position
			ms.pos.y += newY
		} else {
			// collision deteced, set speed to zero so we don't allow the movement
			ms.v.speedY = 0
			if moveY > 0 {
				// reset jump flags
				ms.c.canJump = true
				ms.c.jumping = false
			}
		}

		// subtract from newX and newY the distance we've moved (or not moved)
		moveX += -newX
		moveY += -newY

		// normalize values in case any have gone negative
		if newX > 0 && moveX <= 0 { newX = 0 }
		if newX < 0 && moveX >= 0 { newX = 0 }

		if newY > 0 && moveY <= 0 { newY = 0 }
		if newY < 0 && moveY >= 0 { newY = 0 }

		if moveX == 0 { newX = 0 }
		if moveY == 0 { newY = 0 }

		// if either of these conditions are true, we're done checking
		// and it's time to exit the for loop
		if moveX == 0 && moveY == 0 { break }
		if newX == 0 && newY == 0 { break }
	}
}

func (ms *MovementSystem) StopMove() {
	if ms.v.speedX > 0 {
		ms.v.accelX = -.5
	}

	if ms.v.speedX < 0 {
		ms.v.accelX = .5
	}

	if ms.v.speedX < .2 && ms.v.speedX > -.2 {
		ms.v.accelX = 0
		ms.v.speedX = 0
	}
}

func (ms *MovementSystem) PosValid(newX int32, newY int32) bool {
	TILE_SIZE := ms.engine.Map.TileSize
	retVal := true
	startX := (newX) / TILE_SIZE
	startY := (newY) / TILE_SIZE

	endX := ((newX) + ms.a.w - 1) / TILE_SIZE
	endY := ((newY) + ms.a.h - 1) / TILE_SIZE

	for iY := startY; iY <= endY; iY++ {
		for iX := startX; iX <= endX; iX++ {
			tile := ms.engine.Map.GetTile(iX * TILE_SIZE, iY * TILE_SIZE)

			if ms.PosValidTile(tile) == false {
				retVal = false
			}
		}
	}

	// ENTITY COLLISIONS
	// TODO: Make event system for entity collisions
	/*for i := 0; i < len(EntityList); i++ {
		if ms.PosValidEntity(EntityList[i], newX, newY) == false {
			retVal = false
		}
	}*/

	return retVal
}

func (e *MovementSystem) PosValidTile(tile *Tile) bool {
	if tile == nil {
		return true
	}
	if tile.TypeID == TILE_TYPE_BLOCK {
		return false
	}
	return true
}

/*func (e *MovementSystem) PosValidEntity(gameEntity GameEntity, newX int32, newY int32) bool {
	entity := gameEntity.GetEntity()
	if e != entity &&
		entity != nil &&
		!entity.Dead &&
		entity.Flags & ENTITY_FLAG_MAPONLY == 0 &&
		entity.Collides(newX, newY, e.Image.W - 1, e.Image.H - 1) {
		entityCollision := EntityCollision{entityA: entity, entityB: e}
		EntityCollisionList = append(EntityCollisionList, entityCollision)
		return false
	}

	return true
}*/

type CameraSystem struct {}

func (cs *CameraSystem) Update(engine *Engine, world *World) {
	var pos *Position
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_FOCUSED) {
			pos = &(world.position[entity])
			engine.Camera.SetTarget(&pos.x, &pos.y)
		}
	}
}
