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

func (ss *StateSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_APPEARANCE|COMPONENT_STATE) {
			ap := &(world.appearance[entity])
			stateCmp := &(world.state[entity])
			stateCmp.animationState = stateCmp.animationStates[stateCmp.currentAnimKey]
			if stateCmp.currentAnimKey != stateCmp.desiredAnimKey {
				stateCmp.newState = true
				stateCmp.currentAnimKey = stateCmp.desiredAnimKey
				stateCmp.animationState = stateCmp.animationStates[stateCmp.currentAnimKey]
			}

			ap.name = stateCmp.animationState.asset
			if stateCmp.orientation != stateCmp.animationState.orientation {
				ap.flip = sdl.FLIP_HORIZONTAL
			} else {
				ap.flip = sdl.FLIP_NONE
			}
		}
	}
}

type ColliderSystem struct {}

func (cs *ColliderSystem) Update(engine *Engine, world *World) {
	/*for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_COLLIDER|COMPONENT_APPEARANCE) {
			apCmp := &(world.appearance[entity])
			colCmp := &(world.collider[entity])

			colCmp.
		}
	}*/
}

type RenderSystem struct {
	signature int64
}
func (rs *RenderSystem) Update (engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, RENDER_MASK) {
			pos := &(world.position[entity])
			a := &(world.appearance[entity])
			x := int32(pos.x - engine.Camera.X())
			y := int32((pos.y) - engine.Camera.Y())
			offsetX := x - int32(a.frame.PivotPoint.X / 2)
			offsetY := y - int32(a.frame.PivotPoint.Y / 2)
			if signatureMatches(mask, COMPONENT_COLLIDER) {
				col := &(world.collider[entity])
				engine.Graphics.DrawRectOutline(x, y, col.w, col.h)
			}
			// engine.Graphics.DrawRectOutline(x, y, a.w, a.h)
			engine.Graphics.DrawPart(engine.Assets.Texture, offsetX, offsetY, a.frame.X, a.frame.Y, a.frame.W, a.frame.H, a.flip)
		}
	}
}

type AnimationSystem struct {}

func (as *AnimationSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_ANIMATION|COMPONENT_APPEARANCE|COMPONENT_STATE) {
			animationCmp := &(world.animation[entity])
			appearanceCmp := &(world.appearance[entity])
			stateCmp := &(world.state[entity])

			// determine if we are transitioning to a new state
			// or keeping the current state.  Transitioning to
			// a new state should reset the frame counter to 0
			if stateCmp.newState {
				animationCmp.currentFrame = 0
				stateCmp.newState = false
			}

			// get animation frame array
			frames := engine.Assets.Get(stateCmp.animationState.asset)
			animationCmp.maxFrames = len(frames) - 1
			animationCmp.frameRate = stateCmp.animationState.frameRate
			animationCmp.frameInc = 1
			if len(frames) == 0 { return }

			// calculate offset to be used in rendering to account
			// for different sized sprite frames
			lastFrameIdx := animationCmp.currentFrame - 1
			if lastFrameIdx < 0 {
				lastFrameIdx = len(frames) - 1
			}
			appearanceCmp.frame = frames[animationCmp.currentFrame]

			// advance frames
			if animationCmp.oldTime + uint32(animationCmp.frameRate) > sdl.GetTicks() {
				return
			}
			animationCmp.oldTime = sdl.GetTicks()
			animationCmp.currentFrame += animationCmp.frameInc
			if animationCmp.currentFrame >= animationCmp.maxFrames {
				animationCmp.currentFrame = 0
			}
		}
	}
}

type InputSystem struct {}
func (is *InputSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_STATE) {
			s := &(world.state[entity])

			s.moveRight = engine.Input.KeysHeld[sdl.K_RIGHT] || engine.Input.KeysHeld[sdl.K_d]
			s.moveLeft = engine.Input.KeysHeld[sdl.K_LEFT] || engine.Input.KeysHeld[sdl.K_a]
			s.rolling = (engine.Input.KeysHeld[sdl.K_DOWN] || engine.Input.KeysHeld[sdl.K_s]) && s.grounded
			s.shooting = engine.Input.KeysHeld[sdl.K_RSHIFT]

			if s.moveLeft {
				s.desiredAnimKey = ENTITY_STATE_LEFT
				s.orientation = ORIENTATION_LEFT
			} else if s.moveRight {
				s.desiredAnimKey = ENTITY_STATE_RIGHT
				s.orientation = ORIENTATION_RIGHT
			} else {
				s.desiredAnimKey = ENTITY_STATE_IDLE
			}
			if s.rolling {
				s.desiredAnimKey = ENTITY_STATE_ROLL
			}
			if s.jumping {
				s.desiredAnimKey = ENTITY_STATE_JUMP
			}
			if s.shooting {
				s.desiredAnimKey = ENTITY_STATE_SHOOT
			}
			if engine.Input.KeysHeld[sdl.K_SPACE] {
				s.jumping = true
			} else {
				s.canJump = true
			}
		}
	}
}

type MovementSystem struct {
	v *Velocity
	pos *Position
	stateCmp *State
	a *Appearance
	col *Collider
	engine *Engine
}
func (ms *MovementSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_VELOCITY|COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_STATE|COMPONENT_COLLIDER) {
			ms.v = &(world.velocity[entity])
			ms.pos = &(world.position[entity])
			ms.stateCmp = &(world.state[entity])
			ms.a = &(world.appearance[entity])
			ms.col = &(world.collider[entity])
			ms.engine = engine

			if !ms.stateCmp.moveLeft && !ms.stateCmp.moveRight {
				ms.StopMove()
			}
			if ms.stateCmp.moveLeft {
				ms.v.accelX = -0.2
			} else if ms.stateCmp.moveRight {
				ms.v.accelX = 0.2
			}

			if ms.stateCmp.grounded && ms.stateCmp.jumping && ms.stateCmp.canJump {
				ms.v.speedY = -ms.v.maxSpeedY
				ms.stateCmp.canJump = false
				ms.stateCmp.grounded = false
			}

			if ms.stateCmp.rolling {
				if ms.stateCmp.orientation == ORIENTATION_LEFT {
					ms.v.accelX = -2
				} else if ms.stateCmp.orientation == ORIENTATION_RIGHT {
					ms.v.accelX = 2
				}
			}

			ms.v.accelY = .25
			/*if e.Flags & ENTITY_FLAG_GRAVITY != 0 {
				e.AccelY = .55
			}*/

			ms.v.speedX += ms.v.accelX * engine.FPS.GetSpeedFactor()
			ms.v.speedY += ms.v.accelY * engine.FPS.GetSpeedFactor()

			var maxSpeedX, maxSpeedY float32
			if !ms.stateCmp.rolling {
				maxSpeedX = ms.v.maxSpeedX
				maxSpeedY = ms.v.maxSpeedY
			} else {
				maxSpeedX = ms.v.maxSpeedX + 2
				maxSpeedY = ms.v.maxSpeedY + 2
			}
			if ms.v.speedX > ms.v.maxSpeedX { ms.v.speedX = maxSpeedX }
			if ms.v.speedX < -ms.v.maxSpeedX { ms.v.speedX = -maxSpeedX }
			if ms.v.speedY > ms.v.maxSpeedY { ms.v.speedY = maxSpeedY }
			if ms.v.speedY < -ms.v.maxSpeedY { ms.v.speedY = -maxSpeedY }

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
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.speedX = 0
		}

		// check for collision at desired y
		if ms.PosValid(int32(ms.pos.x), int32(ms.pos.y + newY)) {
			// no collision, update position
			ms.pos.y += newY
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.speedY = 0
			if moveY > 0 {
				// reset jump flags
				ms.stateCmp.grounded = true
				ms.stateCmp.jumping = false
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
		ms.v.accelX = -.3
	}

	if ms.v.speedX < 0 {
		ms.v.accelX = .3
	}

	if ms.v.speedX < .18 && ms.v.speedX > -.18 {
		ms.v.accelX = 0
		ms.v.speedX = 0
	}
}

func (ms *MovementSystem) PosValid(newX int32, newY int32) bool {
	TILE_SIZE := ms.engine.Map.TileSize
	retVal := true
	startX := (newX) / TILE_SIZE
	startY := (newY) / TILE_SIZE

	// endX := ((newX) + ms.a.w - 1) / TILE_SIZE
	// endY := ((newY) + ms.a.h - 1) / TILE_SIZE
	endX := ((newX) + ms.col.w - 1) / TILE_SIZE
	endY := ((newY) + ms.col.h - 1) / TILE_SIZE

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
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_FOCUSED) {
			pos := &(world.position[entity])
			engine.Camera.SetTarget(&pos.x, &pos.y)
		}
	}
}
