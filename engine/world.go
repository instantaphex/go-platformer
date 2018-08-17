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
	controllable [ENTITY_COUNT]Controllable
	focused [ENTITY_COUNT]Focused
	state [ENTITY_COUNT]State
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
			// fmt.Fprintf(os.Stdout, "Entity created: %d\n", entity)
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
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_ANIMATION|COMPONENT_VELOCITY|COMPONENT_CONTROLLABLE|COMPONENT_FOCUSED|COMPONENT_STATE

	w.position[entity].x = x
	w.position[entity].y = y

	w.appearance[entity].name = "Player/Idle"

	w.animation[entity].maxFrames = len(engine.Assets.Get("Player/Idle"))
	w.animation[entity].frameRate = 200
	w.animation[entity].frameInc = 1

	w.velocity[entity].maxSpeedY = 8
	w.velocity[entity].maxSpeedX = 2

	w.controllable[entity].canJump = true

	w.state[entity].currentState = ENTITY_STATE_IDLE
	w.state[entity].stateMap = make(map[EntityStateKey]EntityState)
	w.state[entity].stateMap[ENTITY_STATE_IDLE] = EntityState{ asset: "Player/Idle", flip:  sdl.FLIP_NONE, inheritFlip: true, frameRate: 200 }
	w.state[entity].stateMap[ENTITY_STATE_LEFT] = EntityState{ asset: "Player/Run", flip: sdl.FLIP_HORIZONTAL, inheritFlip: false, frameRate: 60 }
	w.state[entity].stateMap[ENTITY_STATE_RIGHT] = EntityState{ asset: "Player/Run", flip: sdl.FLIP_NONE, inheritFlip: false, frameRate: 60 }
	w.state[entity].stateMap[ENTITY_STATE_JUMP] = EntityState{ asset: "Player/Fall-Jump-WallJ/Jump", flip: sdl.FLIP_NONE, inheritFlip: true, frameRate: 200 }


	return entity
}

func (w *World) CreateStaticPlayer(engine *Engine, x, y float32) int {
	entity := w.CreateEntity()
	w.mask[entity] = COMPONENT_POSITION|COMPONENT_APPEARANCE|COMPONENT_ANIMATION|COMPONENT_VELOCITY

	w.position[entity].x = x
	w.position[entity].y = y

	w.appearance[entity].name = "Player/Idle"

	w.animation[entity].maxFrames = len(engine.Assets.Get("Player/Idle"))
	w.animation[entity].frameRate = 200
	w.animation[entity].frameInc = 1

	w.velocity[entity].maxSpeedY = 8
	w.velocity[entity].maxSpeedX = 3


	return entity
}
