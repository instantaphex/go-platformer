package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type System interface {
	Update(engine *Engine, world *World)
}

func signatureMatches(mask, signature uint64) bool {
	return mask & signature == signature
}

type RenderSystem struct {
}
func (rs *RenderSystem) Update (engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_ANIMATION|COMPONENT_STATE) {
			transformCmp := world.GetTransform(entity)
			animationCmp := world.GetAnimation(entity)
			stateCmp := world.GetState(entity)

			// grab animation metadata
			animState := animationCmp.animationStates[animationCmp.animState]
			// grab frames
			frames := engine.Assets.Get(animState.asset)
			frame := frames[animationCmp.currentFrame]

			// perform sprite flip if needed
			var flip sdl.RendererFlip
			if stateCmp.orientation != animState.orientation {
				flip = sdl.FLIP_HORIZONTAL
			} else {
				flip = sdl.FLIP_NONE
			}

			// determine x and y
			x := int32(transformCmp.x - engine.Camera.X())
			y := int32(transformCmp.y - engine.Camera.Y())

			// line up bounding box center with actual sprite center
			offsetW := (frame.W / 2) - (transformCmp.w / 2)
			offsetH := (frame.H / 2) - (transformCmp.h / 2)
			offsetX := x - offsetW
			offsetY := y - offsetH

			if engine.Config.DrawDebug {
				engine.Graphics.DrawRectOutline(x, y, transformCmp.w, transformCmp.h)
			}
			engine.Graphics.DrawPart(engine.Assets.Texture, offsetX, offsetY, frame.X, frame.Y, frame.W, frame.H, flip)
		}
	}
}

type MapRenderSystem struct {}
func (mrs *MapRenderSystem) Update (engine *Engine, world *World) {
	engine.Map.Render(int32(-engine.Camera.X()), int32(-engine.Camera.Y()))
}

type AnimationSystem struct {}
func (as *AnimationSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_ANIMATION|COMPONENT_STATE) {
			animationCmp := world.GetAnimation(entity)
			stateCmp := world.GetState(entity)

			// determine if we are transitioning to a new state
			// or keeping the current state.  Transitioning to
			// a new state should reset the current frame to 0
			if stateCmp.state != animationCmp.animState {
				animationCmp.currentFrame = 0
				animationCmp.animState = stateCmp.state
			}

			animState := animationCmp.animationStates[animationCmp.animState]
			frames := engine.Assets.Get(animState.asset)
			animationCmp.maxFrames = len(frames)
			animationCmp.frameInc = 1

			// advance frames
			currentTime := sdl.GetTicks()
			threshold := animationCmp.oldTime + uint32(animState.frameRate)
			if threshold > currentTime {
				continue
			}
			animationCmp.oldTime = currentTime
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
		if signatureMatches(mask, COMPONENT_STATE|COMPONENT_CONTROLLER) {
			s := world.GetState(entity)

			s.moveRight = engine.Input.KeysHeld[sdl.K_RIGHT] || engine.Input.KeysHeld[sdl.K_d]
			s.moveLeft = engine.Input.KeysHeld[sdl.K_LEFT] || engine.Input.KeysHeld[sdl.K_a]
			s.rolling = (engine.Input.KeysHeld[sdl.K_DOWN] || engine.Input.KeysHeld[sdl.K_s]) && s.grounded
			s.shooting = engine.Input.KeysHeld[sdl.K_RSHIFT]

			if s.moveLeft {
				s.state = ENTITY_STATE_LEFT
				s.orientation = ORIENTATION_LEFT
			} else if s.moveRight {
				s.state = ENTITY_STATE_RIGHT
				s.orientation = ORIENTATION_RIGHT
			} else {
				s.state = ENTITY_STATE_IDLE
			}
			if s.rolling {
				s.state = ENTITY_STATE_ROLL
			}
			if s.jumping {
				s.state = ENTITY_STATE_JUMP
			}
			if s.shooting {
				s.state = ENTITY_STATE_SHOOT
			}
			if engine.Input.KeysHeld[sdl.K_SPACE] {
				s.jumping = true
			} else {
				s.canJump = true
			}
		}
	}
}

type VelocitySystem struct {
	v *Velocity
	stateCmp *State
}

func (ms *VelocitySystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_VELOCITY|COMPONENT_STATE|COMPONENT_CONTROLLER) {
			ms.v = world.GetVelocity(entity)
			ms.stateCmp = world.GetState(entity)

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

			// apply gravity
			ms.v.accelY = .25

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
		}
	}
}

func (vs *VelocitySystem) StopMove() {
	if vs.v.speedX > 0 {
		vs.v.accelX = -.3
	}

	if vs.v.speedX < 0 {
		vs.v.accelX = .3
	}

	if vs.v.speedX < .18 && vs.v.speedX > -.18 {
		vs.v.accelX = 0
		vs.v.speedX = 0
	}
}

type MovementSystem struct {
	transform *Transform
	v *Velocity
	stateCmp *State
	engine *Engine
	world *World
	currentEntity int
}

func (ms *MovementSystem) Update(engine *Engine, world *World) {
	ms.engine = engine
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_VELOCITY|COMPONENT_STATE) {
			ms.transform = world.GetTransform(entity)
			ms.v = world.GetVelocity(entity)
			ms.stateCmp = world.GetState(entity)
			ms.world = world
			ms.currentEntity = entity
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
		if ms.PosValid(int32(ms.transform.x + newX), int32(ms.transform.y)) {
			// no collision, update position
			ms.transform.x += newX
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.speedX = 0
		}

		// check for collision at desired y
		if ms.PosValid(int32(ms.transform.x), int32(ms.transform.y + newY)) {
			// no collision, update position
			ms.transform.y += newY
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.speedY = 0
			if moveY > 0 {
				// reset jump flags
				// TODO: Possibly raise events for these
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


func (ms *MovementSystem) PosValid(newX int32, newY int32) bool {
	TILE_SIZE := ms.engine.Map.TileSize
	retVal := true
	startX := (newX) / TILE_SIZE
	startY := (newY) / TILE_SIZE

	endX := ((newX) + ms.transform.w - 1) / TILE_SIZE
	endY := ((newY) + ms.transform.h - 1) / TILE_SIZE

	for iY := startY; iY <= endY; iY++ {
		for iX := startX; iX <= endX; iX++ {
			tile := ms.engine.Map.GetTile(iX * TILE_SIZE, iY * TILE_SIZE)

			if ms.PosValidTile(tile) == false {
				retVal = false
			}
		}
	}

	// ENTITY COLLISIONS
	// TODO: Clean this up
	// TODO: Make event system for entity collisions
	collidableEntities := ms.world.GetColliders()
	for _, id := range collidableEntities {
		transformA := ms.world.transform[ms.currentEntity]
		transformB := ms.world.transform[id]
		a := sdl.Rect {
			// check for desired x and y, not current
			X: int32(newX),
			Y: int32(newY),
			W: int32(transformA.w),
			H: int32(transformA.h),
		}
		b := sdl.Rect {
			// check for world coordinates, not camera coordinates
			X: int32(transformB.x - 1),
			Y: int32(transformB.y - 1),
			W: int32(transformB.w),
			H: int32(transformB.h),
		}
		if ms.world.Collides(a, b) && ms.currentEntity != id {
			retVal = false
		}
	}

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

type CameraSystem struct {
	targeted bool
}

func (cs *CameraSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.mask {
		if signatureMatches(mask, COMPONENT_FOCUSED) && !cs.targeted {
			pos := world.GetTransform(entity)
			engine.Camera.SetTarget(&pos.x, &pos.y)
			cs.targeted = true
		}
	}
}
