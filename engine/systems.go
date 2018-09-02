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
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_ANIMATION|COMPONENT_STATE) {
			transformCmp := world.GetTransform(entity)
			animationCmp := world.GetAnimation(entity)
			stateCmp := world.GetState(entity)

			// grab Animation metadata
			animState := animationCmp.AnimationStates[animationCmp.AnimState]
			// grab frames
			frames := engine.Assets.Get(animState.Asset)
			frame := frames[animationCmp.CurrentFrame]

			// perform sprite Flip if needed
			var flip sdl.RendererFlip
			if stateCmp.Orientation != animState.Orientation {
				flip = sdl.FLIP_HORIZONTAL
			} else {
				flip = sdl.FLIP_NONE
			}

			// determine X and Y
			x := int32(transformCmp.X - engine.Camera.X())
			y := int32(transformCmp.Y - engine.Camera.Y())

			// line up bounding box center with actual sprite center
			offsetW := (frame.W / 2) - (transformCmp.W / 2)
			offsetH := (frame.H / 2) - (transformCmp.H / 2)
			offsetX := x - offsetW
			offsetY := y - offsetH

			if engine.Config.DrawDebug {
				engine.Graphics.DrawRectOutline(x, y, transformCmp.W, transformCmp.H)
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
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_ANIMATION|COMPONENT_STATE) {
			animationCmp := world.GetAnimation(entity)
			stateCmp := world.GetState(entity)

			// determine if we are transitioning to a new State
			// or keeping the current State.  Transitioning to
			// a new State should reset the current frame to 0
			if stateCmp.State != animationCmp.AnimState {
				animationCmp.CurrentFrame = 0
				animationCmp.AnimState = stateCmp.State
			}

			animState := animationCmp.CurrentState()
			frames := engine.Assets.Get(animState.Asset)
			animationCmp.MaxFrames = len(frames)
			animationCmp.FrameInc = 1

			// advance frames
			currentTime := sdl.GetTicks()
			threshold := animationCmp.OldTime + uint32(animState.FrameRate)
			if threshold > currentTime {
				continue
			}
			animationCmp.OldTime = currentTime
			animationCmp.CurrentFrame += animationCmp.FrameInc
			if animationCmp.CurrentFrame >= animationCmp.MaxFrames {
				animationCmp.CurrentFrame = 0
			}
		}
	}
}

type InputSystem struct {}
func (is *InputSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_STATE|COMPONENT_CONTROLLER) {
			s := world.GetState(entity)

			s.MoveRight = engine.Input.KeysHeld[sdl.K_RIGHT] || engine.Input.KeysHeld[sdl.K_d]
			s.MoveLeft = engine.Input.KeysHeld[sdl.K_LEFT] || engine.Input.KeysHeld[sdl.K_a]
			s.Rolling = (engine.Input.KeysHeld[sdl.K_DOWN] || engine.Input.KeysHeld[sdl.K_s]) && s.Grounded
			s.Shooting = engine.Input.KeysHeld[sdl.K_RSHIFT]

			if s.MoveLeft {
				s.State = ENTITY_STATE_LEFT
				s.Orientation = ORIENTATION_LEFT
			} else if s.MoveRight {
				s.State = ENTITY_STATE_RIGHT
				s.Orientation = ORIENTATION_RIGHT
			} else {
				s.State = ENTITY_STATE_IDLE
			}
			if s.Rolling {
				s.State = ENTITY_STATE_ROLL
			}
			if s.Jumping {
				s.State = ENTITY_STATE_JUMP
			}
			if s.Shooting {
				s.State = ENTITY_STATE_SHOOT
			}
			if engine.Input.KeysHeld[sdl.K_SPACE] {
				s.Jumping = true
			} else {
				s.CanJump = true
			}
		}
	}
}

type VelocitySystem struct {
	v *Velocity
	stateCmp *State
}

func (ms *VelocitySystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_VELOCITY|COMPONENT_STATE|COMPONENT_CONTROLLER) {
			ms.v = world.GetVelocity(entity)
			ms.stateCmp = world.GetState(entity)

			if !ms.stateCmp.MoveLeft && !ms.stateCmp.MoveRight {
				ms.StopMove()
			}
			if ms.stateCmp.MoveLeft {
				ms.v.AccelX = -0.2
			} else if ms.stateCmp.MoveRight {
				ms.v.AccelX = 0.2
			}

			if ms.stateCmp.Grounded && ms.stateCmp.Jumping && ms.stateCmp.CanJump {
				ms.v.SpeedY = -ms.v.MaxSpeedY
				ms.stateCmp.CanJump = false
				ms.stateCmp.Grounded = false
			}

			if ms.stateCmp.Rolling {
				if ms.stateCmp.Orientation == ORIENTATION_LEFT {
					ms.v.AccelX = -2
				} else if ms.stateCmp.Orientation == ORIENTATION_RIGHT {
					ms.v.AccelX = 2
				}
			}

			// apply gravity
			ms.v.AccelY = .25

			ms.v.SpeedX += ms.v.AccelX * engine.FPS.GetSpeedFactor()
			ms.v.SpeedY += ms.v.AccelY * engine.FPS.GetSpeedFactor()

			var maxSpeedX, maxSpeedY float32
			if !ms.stateCmp.Rolling {
				maxSpeedX = ms.v.MaxSpeedX
				maxSpeedY = ms.v.MaxSpeedY
			} else {
				maxSpeedX = ms.v.MaxSpeedX + 2
				maxSpeedY = ms.v.MaxSpeedY + 2
			}
			if ms.v.SpeedX > ms.v.MaxSpeedX { ms.v.SpeedX = maxSpeedX }
			if ms.v.SpeedX < -ms.v.MaxSpeedX { ms.v.SpeedX = -maxSpeedX }
			if ms.v.SpeedY > ms.v.MaxSpeedY { ms.v.SpeedY = maxSpeedY }
			if ms.v.SpeedY < -ms.v.MaxSpeedY { ms.v.SpeedY = -maxSpeedY }
		}
	}
}

func (vs *VelocitySystem) StopMove() {
	if vs.v.SpeedX > 0 {
		vs.v.AccelX = -.3
	}

	if vs.v.SpeedX < 0 {
		vs.v.AccelX = .3
	}

	if vs.v.SpeedX < .18 && vs.v.SpeedX > -.18 {
		vs.v.AccelX = 0
		vs.v.SpeedX = 0
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
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_VELOCITY|COMPONENT_STATE) {
			ms.transform = world.GetTransform(entity)
			ms.v = world.GetVelocity(entity)
			ms.stateCmp = world.GetState(entity)
			ms.world = world
			ms.currentEntity = entity
			ms.Move(ms.v.SpeedX, ms.v.SpeedY)
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

	// if X or Y are not zero, start collision checking
	// at the current X + speed factor and current Y + speed factor
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

	// starting at speed factor (+ X, + Y) keep iterating over each new desired
	// move, checking for collisions allong the way
	for {
		// check for collision at desired X position
		if ms.PosValid(int32(ms.transform.X+ newX), int32(ms.transform.Y)) {
			// no collision, update position
			ms.transform.X += newX
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.SpeedX = 0
		}

		// check for collision at desired Y
		if ms.PosValid(int32(ms.transform.X), int32(ms.transform.Y+ newY)) {
			// no collision, update position
			ms.transform.Y += newY
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.v.SpeedY = 0
			if moveY > 0 {
				// reset jump flags
				// TODO: Possibly raise events for these
				ms.stateCmp.Grounded = true
				ms.stateCmp.Jumping = false
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

	endX := ((newX) + ms.transform.W - 1) / TILE_SIZE
	endY := ((newY) + ms.transform.H - 1) / TILE_SIZE

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
		transformA := ms.world.GetTransform(ms.currentEntity)
		transformB := ms.world.GetTransform(id)
		a := transformA.GetPotentialBB(newX, newY)
		b := transformB.GetBB()
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
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_FOCUSED) && !cs.targeted {
			pos := world.GetTransform(entity)
			engine.Camera.SetTarget(&pos.X, &pos.Y)
			cs.targeted = true
		}
	}
}
