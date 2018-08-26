package engine

import (
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
)

const ENTITY_COUNT = 5

type World struct {
	mask [ENTITY_COUNT]uint64
	position [ENTITY_COUNT]Position
	velocity [ENTITY_COUNT]Velocity
	appearance [ENTITY_COUNT]Appearance
	animation [ENTITY_COUNT]Animation
	focused [ENTITY_COUNT]Focused
	state [ENTITY_COUNT]State
	collider [ENTITY_COUNT]Collider
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

func (w *World) CreatePlayer(engine *Engine, x, y float32) int {
	entity := w.CreateEntity()
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_ANIMATION|COMPONENT_VELOCITY|COMPONENT_FOCUSED|COMPONENT_STATE|COMPONENT_COLLIDER

	w.position[entity].x = x
	w.position[entity].y = y

	w.appearance[entity].name = "Player/Idle"

	w.animation[entity].maxFrames = len(engine.Assets.Get("Player/Idle"))
	w.animation[entity].frameRate = 200
	w.animation[entity].frameInc = 1

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
