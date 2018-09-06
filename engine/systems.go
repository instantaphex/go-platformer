package engine

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"strconv"
)

type System interface {
	Update(engine *Engine, world *World)
	Init(world *World)
}

type SystemEvents struct {
	events []Event
}

func (se *SystemEvents) PostEvent(e Event) {
	se.events = append(se.events, e)
}

func (se *SystemEvents) HandleEvents(handler func(Event)) {
	for _, event := range se.events {
		handler(event)
	}
	se.events = nil
}

func signatureMatches(mask, signature uint64) bool {
	return mask & signature == signature
}

type RenderSystem struct {
	SystemEvents
}
func (rs *RenderSystem) Init(world *World) {
}
func (rs *RenderSystem) Update (engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_ANIMATION) {
			transformCmp := world.GetTransform(entity)
			animationCmp := world.GetAnimation(entity)
			stateCmp := world.GetState(entity)

			// grab Animation metadata
			animState := animationCmp.AnimationStates[animationCmp.AnimState]
			// grab frames
			frames := engine.Assets.Get(animState.Asset)
			frame := frames[animationCmp.CurrentFrame]

			// perform sprite Flip if needed
			flip := sdl.FLIP_NONE
			if signatureMatches(mask, COMPONENT_STATE) {
				if stateCmp.Orientation != animState.Orientation {
					flip = sdl.FLIP_HORIZONTAL
				} else {
					flip = sdl.FLIP_NONE
				}
			}

			// determine X and Y
			x := int32(transformCmp.X - engine.Camera.X())
			y := int32(transformCmp.Y - engine.Camera.Y())

			if signatureMatches(mask, COMPONENT_HUD) {
				x = int32(transformCmp.X)
				y = int32(transformCmp.Y)
			}

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

type MapRenderSystem struct {
	SystemEvents
}
func (mrs *MapRenderSystem) Init(world *World) {

}
func (mrs *MapRenderSystem) Update (engine *Engine, world *World) {
	engine.Map.Render(int32(-engine.Camera.X()), int32(-engine.Camera.Y()))
}

type AnimationSystem struct {
	SystemEvents
}
func (as *AnimationSystem) Init(world *World) {

}
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

type InputSystem struct {
	SystemEvents
}
func (is *InputSystem) Init(world *World) {}
func (is *InputSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_STATE|COMPONENT_CONTROLLER|COMPONENT_TRANSFORM) {
			s := world.GetState(entity)
			transform := world.GetTransform(entity)

			// s.MoveRight = engine.Input.KeysHeld[sdl.K_RIGHT] || engine.Input.KeysHeld[sdl.K_d]
			s.MoveRight = engine.Input.KeysHeld[sdl.K_RIGHT] && !s.Sliding
			s.MoveLeft = engine.Input.KeysHeld[sdl.K_LEFT] && !s.Sliding
			s.Rolling = engine.Input.KeysHeld[sdl.K_DOWN] && s.Grounded
			s.Shooting = engine.Input.KeysHeld[sdl.K_RSHIFT]

			s.Grounded = transform.Sensor.Bottom
			s.LeftSlide = transform.Sensor.Left
			s.RightSlide = transform.Sensor.Right
			s.Sliding = (s.LeftSlide || s.RightSlide) && !s.Grounded

			// s.Jumping = !s.Grounded

			if s.Grounded {
				s.Jumping = false
				s.JumpCount = 0
				s.JumpFrameCount = 0
			}

			if s.MoveLeft {
				s.State = ENTITY_STATE_LEFT
				s.Orientation = ORIENTATION_LEFT
			} else if s.MoveRight  {
				s.State = ENTITY_STATE_RIGHT
				s.Orientation = ORIENTATION_RIGHT
			} else {
				s.State = ENTITY_STATE_IDLE
			}

			if !s.Grounded {
				s.State = ENTITY_STATE_JUMP
			}
			if s.Rolling {
				s.State = ENTITY_STATE_ROLL
			}
			if s.Shooting {
				s.State = ENTITY_STATE_SHOOT
			}

			/**
			 * Wall Sliding
			 */
			if s.RightSlide && !s.Grounded {
				s.State = ENTITY_STATE_WALLR
				s.JumpCount = 0
			}
			if s.LeftSlide && !s.Grounded {
				s.State = ENTITY_STATE_WALLL
				s.JumpCount = 0
			}

			/**
			 * Jumping
			 */
			if engine.Input.KeysHeld[sdl.K_SPACE] {
				s.JumpFrameCount++
			}

			if engine.Input.KeyState(sdl.K_SPACE).JustPressed()  && s.JumpCount < 2 {
				s.Jumping = true

				// For jumping purposes sliding is the same as being on the ground
				if !s.Sliding {
					s.JumpCount++
				}

				// determine jump height and
				var speedX, speedY float32
				speedY = -5
				if s.LeftSlide && !s.Grounded {
					speedX = 70
					s.State = ENTITY_STATE_RIGHT
					s.Orientation = ORIENTATION_RIGHT
				} else if s.RightSlide && !s.Grounded {
					speedX = -5
					s.State = ENTITY_STATE_LEFT
					s.Orientation = ORIENTATION_LEFT
				} else {
					speedX = 0
				}

				// emit physics pulse event
				world.Events.EmitEvent(&PhysicsPulseEvent{
					Entity: entity,
					SpeedX: speedX,
					SpeedY: speedY / float32(s.JumpCount / 2),
				})
				// Emit jump Audio Event
				world.Events.EmitEvent(&AudioEvent{ Clip: "jump.wav" })
			}

			if s.Jumping && engine.Input.KeyState(sdl.K_SPACE).JustReleased() && !s.Sliding && s.JumpFrameCount < 25 {
				s.Jumping = false
				s.JumpFrameCount = 0
				world.Events.EmitEvent(&PhysicsPulseEvent{
					Entity: entity,
					SpeedX: 0,
					SpeedY: 2,
				})
			}
		}
	}
}

type PhysicsSystem struct {
	transform *Transform
	stateCmp *State
	SystemEvents
}

func (ps *PhysicsSystem) Init(world *World) {
	world.Events.Subscribe("physics-pulse", ps)
}
func (ps *PhysicsSystem) Update(engine *Engine, world *World) {
	ps.SystemEvents.HandleEvents(func(event Event) {
		evt, _ := event.(*PhysicsPulseEvent)
		mask := world.Mask[evt.Entity]
		if signatureMatches(mask, COMPONENT_VELOCITY|COMPONENT_STATE|COMPONENT_CONTROLLER) {
			transform := world.GetTransform(evt.Entity)
			transform.SpeedX += evt.SpeedX
			transform.SpeedY += evt.SpeedY
		}
	})
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_STATE|COMPONENT_CONTROLLER) {
			ps.transform = world.GetTransform(entity)
			ps.stateCmp = world.GetState(entity)

			if !ps.stateCmp.MoveLeft && !ps.stateCmp.MoveRight {
				ps.StopMove()
			}
			if ps.stateCmp.MoveLeft {
				ps.transform.AccelX = -0.2
			} else if ps.stateCmp.MoveRight {
				ps.transform.AccelX = 0.2
			}

			if ps.stateCmp.Rolling {
				if ps.stateCmp.Orientation == ORIENTATION_LEFT {
					ps.transform.AccelX = -2
				} else if ps.stateCmp.Orientation == ORIENTATION_RIGHT {
					ps.transform.AccelX = 2
				}
			}

			// apply gravity
			ps.transform.AccelY = .13

			if ps.stateCmp.Sliding && ps.transform.SpeedY > 0 {
				ps.transform.SpeedY /= 2
			}

			ps.transform.SpeedX += ps.transform.AccelX * engine.FPS.GetSpeedFactor()
			ps.transform.SpeedY += ps.transform.AccelY * engine.FPS.GetSpeedFactor()

			var maxSpeedX, maxSpeedY float32
			if !ps.stateCmp.Rolling {
				maxSpeedX = ps.transform.MaxSpeedX
				maxSpeedY = ps.transform.MaxSpeedY
			} else {
				maxSpeedX = ps.transform.MaxSpeedX + 2
				maxSpeedY = ps.transform.MaxSpeedY + 2
			}
			if ps.transform.SpeedX > ps.transform.MaxSpeedX { ps.transform.SpeedX = maxSpeedX }
			if ps.transform.SpeedX < -ps.transform.MaxSpeedX { ps.transform.SpeedX = -maxSpeedX }
			if ps.transform.SpeedY > ps.transform.MaxSpeedY { ps.transform.SpeedY = maxSpeedY }
			if ps.transform.SpeedY < -ps.transform.MaxSpeedY { ps.transform.SpeedY = -maxSpeedY }
		}
	}
}

func (ps *PhysicsSystem) StopMove() {
	if ps.transform.SpeedX > 0 {
		ps.transform.AccelX = -.3
	}

	if ps.transform.SpeedX < 0 {
		ps.transform.AccelX = .3
	}

	if ps.transform.SpeedX < .18 && ps.transform.SpeedX > -.18 {
		ps.transform.AccelX = 0
		ps.transform.SpeedX = 0
	}
}

type MovementSystem struct {
	transform *Transform
	stateCmp *State
	engine *Engine
	world *World
	currentEntity int
	SystemEvents
}
func (ms *MovementSystem) Init(world *World) {}
func (ms *MovementSystem) Update(engine *Engine, world *World) {
	ms.engine = engine
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_VELOCITY|COMPONENT_STATE) {
			ms.transform = world.GetTransform(entity)
			ms.stateCmp = world.GetState(entity)
			ms.world = world
			ms.currentEntity = entity
			ms.Move(ms.transform.SpeedX, ms.transform.SpeedY)
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
	// move, checking for collisions along the way
	for {
		// check for collision at desired X position
		if ms.PosValid(int32(ms.transform.X+ newX), int32(ms.transform.Y)) {
			// no collision, update position
			ms.transform.X += newX
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.transform.SpeedX = 0
		}

		// check for collision at desired Y
		if ms.PosValid(int32(ms.transform.X), int32(ms.transform.Y+ newY)) {
			// no collision, update position
			ms.transform.Y += newY
		} else {
			// collision detected, set speed to zero so we don't allow the movement
			ms.transform.SpeedY = 0
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
	retVal := true
	startX := newX
	startY := newY

	endX := newX + ms.transform.W - 1
	endY := newY + ms.transform.H - 1

	for iY := startY; iY <= endY; iY++ {
		for iX := startX; iX <= endX; iX++ {
			if ms.engine.Map.PointCollidesTile(iX, iY) {
				retVal = false
			}
		}
	}

	// sensors
	// TODO: Make bottom sensor (at least) a line that extends full width of the entity BB
	top, bottom, left, right := ms.transform.GetSensorPoints()
	ms.transform.Sensor.Top = ms.engine.Map.PointCollidesTile(top.X, top.Y)
	ms.transform.Sensor.Bottom = ms.engine.Map.PointCollidesTile(bottom.X, bottom.Y)
	ms.transform.Sensor.Left = ms.engine.Map.PointCollidesTile(left.X, left.Y)
	ms.transform.Sensor.Right = ms.engine.Map.PointCollidesTile(right.X, right.Y)

	return retVal
}

type CameraSystem struct {
	targeted bool
	SystemEvents
}
func (cs *CameraSystem) Init(world *World) {}
func (cs *CameraSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_FOCUSED) && !cs.targeted {
			pos := world.GetTransform(entity)
			engine.Camera.SetTarget(&pos.X, &pos.Y)
			cs.targeted = true
		}
	}
}

type EntityCollisionSystem struct {
	SystemEvents
}
func (ecs *EntityCollisionSystem) Init(world *World) {}
func (ecs *EntityCollisionSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_CONTROLLER|COMPONENT_TRANSFORM) {
			collidableEntities := world.GetColliders()
			for _, id := range collidableEntities {
				transformA := world.GetTransform(entity)
				transformB := world.GetTransform(id)
				a := transformA.GetBB()
				b := transformB.GetBB()
				if world.Collides(a, b) && entity != id {
					world.Events.EmitEvent(&CollisionEvent{ A: entity, B: id })
				}
			}
		}
	}
}

type EntityCollectionSystem struct {
	SystemEvents
	subs []Subscription
}
func (ecs *EntityCollectionSystem) Init(world *World) {
	ecs.subs = append(ecs.subs, world.Events.Subscribe("collision", ecs))
}
func (ecs *EntityCollectionSystem) Update(engine *Engine, world *World) {
	ecs.SystemEvents.HandleEvents(func(event Event) {
		evt, ok := event.(*CollisionEvent)
		if !ok {
			fmt.Fprintf(os.Stderr, "Event type: %s is not a valid Collsision Event", event.Type())
		}
		// fmt.Fprintf(os.Stdout, "Collision event A: %d\t\t\tB:%d\n", evt.A, evt.B)
		maskA := world.Mask[evt.A]
		maskB := world.Mask[evt.B]

		if signatureMatches(maskA, COMPONENT_INVENTORY) && signatureMatches(maskB, COMPONENT_COLLECTIBLE) {
			world.DestroyEntity(evt.B)
			inventory := world.Inventory[evt.A]
			collectible := world.Collectible[evt.B]
			inventory.Items[collectible.Type] += collectible.Value
			world.Events.EmitEvent(&CollectionEvent{
				Collectible: collectible.Type,
				Total: inventory.Items[collectible.Type],
				NumCollected: collectible.Value,
			})
		}
	})
}

type AudioSystem struct {
	SystemEvents
	sub Subscription
}
func (as *AudioSystem) Init(world *World) {
	as.sub = world.Events.Subscribe("collection", as)
	as.sub = world.Events.Subscribe("audio", as)
}
func (as *AudioSystem) Update(engine *Engine, world *World) {
	as.HandleEvents(func(event Event) {
		switch evt := event.(type) {
		case *CollectionEvent:
			if evt.Collectible == "gold" {
				engine.Audio.PlaySoundEffect("coin.wav")
			}
			if evt.Collectible == "health" {
				engine.Audio.PlaySoundEffect("health.wav")
			}
		case *AudioEvent:
			engine.Audio.PlaySoundEffect(evt.Clip)
		default:

		}
	})
}

type TextRenderSystem struct {
	SystemEvents
}
func (trs *TextRenderSystem) Init(world *World) {}
func (trs *TextRenderSystem) Update(engine *Engine, world *World) {
	for entity, mask := range world.Mask {
		if signatureMatches(mask, COMPONENT_TRANSFORM|COMPONENT_TEXT) {
			transform := world.GetTransform(entity)
			text := world.GetText(entity)
			cacheId := strconv.Itoa(entity)
			engine.Text.Write(cacheId, text.Value)
			texture := engine.Text.GetTexture(cacheId)
			if texture != nil {
				engine.Graphics.DrawFull(texture, int32(transform.X), int32(transform.Y))
			}
		}
	}
}

type HudTextSystem struct {
	SystemEvents
}
func (chs *HudTextSystem) Init(world *World) {
	world.Events.Subscribe("collection", chs)
}
func (chs *HudTextSystem) Update(engine *Engine, world *World) {
	chs.SystemEvents.HandleEvents(func (event Event) {
		evt, _ := event.(*CollectionEvent);
		if evt.Collectible == "gold" {
			text := world.GetTextByTag("player_coins")
			text.Value = "Coins: " + strconv.Itoa(evt.Total)
		}
		if evt.Collectible == "health" {
			text := world.GetTextByTag("player_health")
			text.Value = "X " + strconv.Itoa(evt.Total)
		}
	})
}

