package engine

import (
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
)

const ENTITY_COUNT = 3

type World struct {
	mask [ENTITY_COUNT]uint64
	position [ENTITY_COUNT]Position
	velocity [ENTITY_COUNT]Velocity
	appearance [ENTITY_COUNT]Appearance
	animation [ENTITY_COUNT]Animation
	focused [ENTITY_COUNT]Focused
	state [ENTITY_COUNT]State
	collider [ENTITY_COUNT]Collider
	controller [ENTITY_COUNT]Controller
	systems []System
}

func (w *World) RegisterSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) Update(engine *Engine) {
	for _, system := range w.systems {
		system.Update(engine, w)
	}
}

func (w *World) CreateEntity() int {
	for entity := 0; entity < ENTITY_COUNT; entity++ {
		if w.mask[entity] == COMPONENT_NONE {
			return entity
		}
	}

	fmt.Fprintf(os.Stderr, "No more entities left\n")
	return ENTITY_COUNT
}

func (w *World) DestroyEntity(entity uint64) {
	fmt.Fprintf(os.Stdout, "Entity destroyed: %d\n", entity)
	w.mask[entity] = COMPONENT_NONE
}

func (w *World) GetColliders() []int {
	var list []int
	for entity, signature := range w.mask {
		if signatureMatches(signature, COMPONENT_COLLIDER) {
			list = append(list, entity)
		}
	}
	return list
}

func (w *World) Collides(a, b sdl.Rect) bool {
	var left1, left2, right1, right2, top1, top2, bottom1, bottom2 int32

	left1 = a.X
	left2 = b.X

	right1 = left1 + a.W - 1
	right2 = b.X + b.W - 1

	top1 = a.Y
	top2 = b.Y

	bottom1 = top1 + a.H - 1
	bottom2 = top2 + b.H - 1

	if bottom1 < top2 { return false }
	if top1 > bottom2 { return false }
	if right1 < left2 { return false }
	if left1 > right2 { return false }

	return true
}

func (w *World) CreateHeart(engine * Engine, x, y float32) int {
	entity := w.CreateEntity()
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_COLLIDER|COMPONENT_ANIMATION|COMPONENT_STATE
	w.position[entity].x = x
	w.position[entity].y = y
	w.collider[entity].w = 8
	w.collider[entity].h = 7
	w.state[entity].animationStates = make(map[AnimationStateKey]AnimationState)
	w.state[entity].animationStates[ENTITY_STATE_IDLE] = AnimationState{
		asset: "Items/Heart/Pick heart",
		flip:  sdl.FLIP_NONE,
		frameRate: 200,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}
	return entity
}

func (w *World) CreateCoin(engine * Engine, x, y float32) int {
	entity := w.CreateEntity()
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_COLLIDER|COMPONENT_ANIMATION|COMPONENT_STATE
	w.position[entity].x = x
	w.position[entity].y = y
	w.collider[entity].w = 8
	w.collider[entity].h = 8
	w.state[entity].animationStates = make(map[AnimationStateKey]AnimationState)
	w.state[entity].animationStates[ENTITY_STATE_IDLE] = AnimationState{
		asset: "Items/Coin/Shine",
		flip:  sdl.FLIP_NONE,
		frameRate: 200,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}
	return entity
}

func (w *World) CreatePlayer(engine *Engine, x, y float32) int {
	entity := w.CreateEntity()
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_ANIMATION|COMPONENT_VELOCITY|COMPONENT_FOCUSED|COMPONENT_STATE|COMPONENT_COLLIDER|COMPONENT_CONTROLLER

	w.position[entity].x = x
	w.position[entity].y = y

	w.velocity[entity].maxSpeedY = 5
	w.velocity[entity].maxSpeedX = 2.2

	w.collider[entity].w = 9
	w.collider[entity].h = 14

	w.state[entity].canJump = true

	w.state[entity].currentAnimKey = ENTITY_STATE_IDLE

	w.state[entity].animationStates = make(map[AnimationStateKey]AnimationState)
	w.state[entity].animationStates[ENTITY_STATE_IDLE] = AnimationState{
		asset: "Player/Idle",
		flip:  sdl.FLIP_NONE,
		frameRate: 200,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}
	w.state[entity].animationStates[ENTITY_STATE_LEFT] = AnimationState{
		asset: "Player/Run",
		flip: sdl.FLIP_HORIZONTAL,
		frameRate: 60,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}
	w.state[entity].animationStates[ENTITY_STATE_RIGHT] = AnimationState{
		asset: "Player/Run",
		flip: sdl.FLIP_NONE,
		frameRate: 60,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}
	w.state[entity].animationStates[ENTITY_STATE_JUMP] = AnimationState{
		asset: "Player/Fall-Jump-WallJ/Jump",
		flip: sdl.FLIP_NONE,
		frameRate: 0,
		infinite: true,
		orientation: ORIENTATION_RIGHT,
	}

	w.state[entity].animationStates[ENTITY_STATE_ROLL] = AnimationState{
		asset: "Player/Roll",
		flip: sdl.FLIP_NONE,
		frameRate: 150,
		infinite: false,
		orientation: ORIENTATION_RIGHT,
	}

	w.state[entity].animationStates[ENTITY_STATE_SHOOT] = AnimationState{
		asset: "Player/Bow",
		flip: sdl.FLIP_NONE,
		frameRate: 150,
		infinite: false,
		orientation: ORIENTATION_RIGHT,
	}

	return entity
}
