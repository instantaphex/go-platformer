package engine

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const ENTITY_COUNT = 100

type EntityBuilder func(world *World, x, y float32) int

type World struct {
	// data components
	Mask        	[ENTITY_COUNT]uint64
	Transform   	[ENTITY_COUNT]Transform
	Animation   	[ENTITY_COUNT]Animation
	State       	[ENTITY_COUNT]State
	Tag		    	[ENTITY_COUNT]Tag
	Collectible 	[ENTITY_COUNT]Collectible
	Inventory   	[ENTITY_COUNT]Inventory
	Text        	[ENTITY_COUNT]Text

	// no data components
	Focused     	[ENTITY_COUNT]Focused
	Controller  	[ENTITY_COUNT]Controller
	Collidable  	[ENTITY_COUNT]Collidable
	Hud         	[ENTITY_COUNT]Hud

	systems []System
	entityBuilders map[string]EntityBuilder
	Events *Dispatcher
}

func NewWorld() *World {
	return &World {
		Events: NewDispatcher(),
	}
}

func (w *World) RegisterSystem(system System) {
	if w.systems == nil {
		w.systems = make([]System, 0)
	}
	system.Init(w)
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

func (w *World) GetAnimation(entity int) *Animation {
	return &w.Animation[entity]
}

func (w *World) GetState(entity int) *State {
	return &w.State[entity]
}

func (w *World) GetTag(entity int) *Tag {
	return &w.Tag[entity]
}

func (w *World) GetCollectible(entity int) *Collectible {
	return &w.Collectible[entity]
}

func (w *World) GetInventory(entity int) *Inventory {
	return &w.Inventory[entity]
}

func (w *World) GetText(entity int) *Text {
	return &w.Text[entity]
}

func (w *World) GetTextByTag(value string) *Text {
	for entity, mask := range w.Mask {
		if signatureMatches(mask, COMPONENT_TAG|COMPONENT_TEXT) {
			tag := w.GetTag(entity)
			if tag.Value == value {
				return w.GetText(entity)
			}
		}
	}
	return nil
}

func (w *World) GetHud(entity int) *Hud {
	return &w.Hud[entity]
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

func (w *World) DestroyEntity(entity int) {
	// fmt.Fprintf(os.Stdout, "Entity destroyed: %d\n", entity)
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
