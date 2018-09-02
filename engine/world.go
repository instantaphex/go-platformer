package engine

import (
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
)

const ENTITY_COUNT = 3

type EntityBuilder func(world *World, x, y float32) int

type World struct {
	// data components
	Mask      [ENTITY_COUNT]uint64
	Transform [ENTITY_COUNT]Transform
	Velocity  [ENTITY_COUNT]Velocity
	Animation [ENTITY_COUNT]Animation
	State     [ENTITY_COUNT]State

	// no data components
	Focused    [ENTITY_COUNT]Focused
	Controller [ENTITY_COUNT]Controller
	Collidable [ENTITY_COUNT]Collidable

	systems []System
	entityBuilders map[string]EntityBuilder
}

func (w *World) RegisterSystem(system System) {
	if w.systems == nil {
		w.systems = make([]System, 0)
	}
	w.systems = append(w.systems, system)
}

func (w *World) RegisterEntityBuilder(name string, builder EntityBuilder) {
	if w.entityBuilders == nil {
		w.entityBuilders = make(map[string]EntityBuilder)
	}
	w.entityBuilders[name] = builder
}

func (w *World) Update(engine *Engine) {
	for _, system := range w.systems {
		system.Update(engine, w)
	}
}

func (w *World) GetMask(entity int) *uint64 {
	return &w.Mask[entity]
}

func (w *World) GetTransform(entity int) *Transform {
	return &w.Transform[entity]
}

func (w *World) GetVelocity(entity int) *Velocity {
	return &w.Velocity[entity]
}

func (w *World) GetAnimation(entity int) *Animation {
	return &w.Animation[entity]
}

func (w *World) GetState(entity int) *State {
	return &w.State[entity]
}

func (w *World) CreateEntity() int {
	for entity := 0; entity < ENTITY_COUNT; entity++ {
		if w.Mask[entity] == COMPONENT_NONE {
			return entity
		}
	}

	fmt.Fprintf(os.Stderr, "No more entities left\n")
	return ENTITY_COUNT
}

func (w *World) DestroyEntity(entity uint64) {
	fmt.Fprintf(os.Stdout, "Entity destroyed: %d\n", entity)
	w.Mask[entity] = COMPONENT_NONE
}

func (w *World) GetColliders() []int {
	var list []int
	for entity, signature := range w.Mask {
		if signatureMatches(signature, COMPONENT_TRANSFORM) {
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
