package main

import (
	"github.com/instantaphex/platformer/engine"
	"github.com/veandco/go-sdl2/sdl"
)

func  CreateHeart(w *engine.World, x, y float32) int {
	entity := w.CreateEntity()
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_STATE
	w.Transform[entity].X = x
	w.Transform[entity].Y = y
	w.Transform[entity].W = 8
	w.Transform[entity].H = 7
	w.Animation[entity].AnimationStates = make(map[engine.StateKey]engine.AnimationState)
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_IDLE] = engine.AnimationState{
		Asset: "Items/Heart/Pick heart",
		Flip:  sdl.FLIP_NONE,
		FrameRate: 100,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	return entity
}

func CreateCoin(w *engine.World, x, y float32) int {
	entity := w.CreateEntity()
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_STATE
	w.Transform[entity].X = x
	w.Transform[entity].Y = y
	w.Transform[entity].W = 8
	w.Transform[entity].H = 8
	w.Animation[entity].AnimationStates = make(map[engine.StateKey]engine.AnimationState)
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_IDLE] = engine.AnimationState{
		Asset: "Items/Coin/Shine",
		Flip:  sdl.FLIP_NONE,
		FrameRate: 200,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	return entity
}

func CreatePlayer(w *engine.World, x, y float32) int {
	entity := w.CreateEntity()
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_VELOCITY|engine.COMPONENT_FOCUSED|engine.COMPONENT_STATE|engine.COMPONENT_CONTROLLER

	w.Transform[entity].X = x
	w.Transform[entity].Y = y
	w.Transform[entity].W = 9
	w.Transform[entity].H = 14

	w.Velocity[entity].MaxSpeedY = 5
	w.Velocity[entity].MaxSpeedX = 2.2

	w.State[entity].CanJump = true

	w.State[entity].State = engine.ENTITY_STATE_IDLE

	w.Animation[entity].AnimationStates = make(map[engine.StateKey]engine.AnimationState)
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_IDLE] = engine.AnimationState{
		Asset: "Player/Idle",
		Flip:  sdl.FLIP_NONE,
		FrameRate: 200,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_LEFT] = engine.AnimationState{
		Asset: "Player/Run",
		Flip: sdl.FLIP_HORIZONTAL,
		FrameRate: 60,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_RIGHT] = engine.AnimationState{
		Asset: "Player/Run",
		Flip: sdl.FLIP_NONE,
		FrameRate: 60,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_JUMP] = engine.AnimationState{
		Asset: "Player/Fall-Jump-WallJ/Jump",
		Flip: sdl.FLIP_NONE,
		FrameRate: 0,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}

	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_ROLL] = engine.AnimationState{
		Asset: "Player/Roll",
		Flip: sdl.FLIP_NONE,
		FrameRate: 150,
		Infinite: false,
		Orientation: engine.ORIENTATION_RIGHT,
	}

	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_SHOOT] = engine.AnimationState{
		Asset: "Player/Bow",
		Flip: sdl.FLIP_NONE,
		FrameRate: 150,
		Infinite: false,
		Orientation: engine.ORIENTATION_RIGHT,
	}

	return entity
}

func main() {
	eng := engine.New(engine.EngineConfig{
		WindowWidth: 1024,
		WindowHeight: 768,
		WindowTitle: "Platformer",
		Scale: 2,
		DrawDebug: false,
	})

	eng.World = &engine.World{}
	eng.World.RegisterSystem(&engine.CameraSystem{})
	eng.World.RegisterSystem(&engine.InputSystem{})
	eng.World.RegisterSystem(&engine.AnimationSystem{})
	eng.World.RegisterSystem(&engine.VelocitySystem{})
	eng.World.RegisterSystem(&engine.MovementSystem{})

	// order matters here
	eng.World.RegisterSystem(&engine.MapRenderSystem{})
	eng.World.RegisterSystem(&engine.RenderSystem{})

	eng.World.RegisterEntityBuilder("coin", CreateCoin)
	eng.World.RegisterEntityBuilder("player", CreatePlayer)
	eng.World.RegisterEntityBuilder("heart", CreateHeart)

	eng.Map.Load("level1", eng.World)

	eng.Run()
}
