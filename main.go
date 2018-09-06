package main

import (
	"github.com/instantaphex/platformer/engine"
	"github.com/veandco/go-sdl2/sdl"
)

func  CreateHeart(w *engine.World, x, y float32) int {
	entity := w.CreateEntity()
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_STATE|engine.COMPONENT_TAG|engine.COMPONENT_COLLECTIBLE
	w.Tag[entity].Value = "heart"
	w.Collectible[entity].Type = "health"
	w.Collectible[entity].Value = 1
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
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_STATE|engine.COMPONENT_TAG|engine.COMPONENT_COLLECTIBLE
	w.Tag[entity].Value = "coin"
	w.Collectible[entity].Type = "gold"
	w.Collectible[entity].Value = 1
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
	w.Mask[entity] = engine.COMPONENT_TRANSFORM|engine.COMPONENT_ANIMATION|engine.COMPONENT_VELOCITY|engine.COMPONENT_FOCUSED|engine.COMPONENT_STATE|engine.COMPONENT_CONTROLLER|engine.COMPONENT_TAG|engine.COMPONENT_INVENTORY

	w.Tag[entity].Value = "player"

	w.Inventory[entity].Items = make(map[string]int)
	w.Inventory[entity].Items["health"] = 3
	w.Inventory[entity].Items["gold"] = 0

	w.Transform[entity].X = x
	w.Transform[entity].Y = y
	w.Transform[entity].W = 9
	w.Transform[entity].H = 14

	w.Transform[entity].MaxSpeedY = 4
	w.Transform[entity].MaxSpeedX = 2.2

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

	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_WALLR] = engine.AnimationState{
		Asset: "Player/Fall-Jump-WallJ/WallJ",
		Flip: sdl.FLIP_HORIZONTAL,
		FrameRate: 0,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}

	w.Animation[entity].AnimationStates[engine.ENTITY_STATE_WALLL] = engine.AnimationState{
		Asset: "Player/Fall-Jump-WallJ/WallJ",
		Flip: sdl.FLIP_HORIZONTAL,
		FrameRate: 0,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}

	return entity
}

func CreateScoreHud(w *engine.World) {
	entity := w.CreateEntity()

	w.Mask[entity] = engine.COMPONENT_TEXT|engine.COMPONENT_TRANSFORM|engine.COMPONENT_TAG

	w.Text[entity].Value = "Coins: 0"

	w.Tag[entity].Value = "player_coins"

	w.Transform[entity].X = 75
	w.Transform[entity].Y = 2
}

func CreateHealthHud(w *engine.World) {
	heart := w.CreateEntity()
	w.Mask[heart] = engine.COMPONENT_ANIMATION|engine.COMPONENT_TRANSFORM|engine.COMPONENT_STATE|engine.COMPONENT_HUD
	w.Animation[heart].AnimationStates = make(map[engine.StateKey]engine.AnimationState)
	w.Animation[heart].AnimationStates[engine.ENTITY_STATE_IDLE] = engine.AnimationState {
		Asset: "Items/Heart/heart-red",
		Flip:  sdl.FLIP_NONE,
		FrameRate: 0,
		Infinite: true,
		Orientation: engine.ORIENTATION_RIGHT,
	}
	w.Transform[heart].X = 0
	w.Transform[heart].Y = 0
	w.Transform[heart].W = 10
	w.Transform[heart].H = 10
	w.State[heart].State = engine.ENTITY_STATE_IDLE

	num := w.CreateEntity()
	w.Mask[num] = engine.COMPONENT_TEXT|engine.COMPONENT_TRANSFORM|engine.COMPONENT_TAG
	w.Text[num].Value = "X 3"
	w.Tag[num].Value = "player_health"
	w.Transform[num].X = 18
	w.Transform[num].Y = 2
}

func main() {
	eng := engine.New(engine.EngineConfig{
		WindowWidth: 1024,
		WindowHeight: 768,
		WindowTitle: "Platformer",
		Scale: 2,
		DrawDebug: false,
	})

	eng.World = engine.NewWorld()
	eng.World.RegisterSystem(&engine.CameraSystem{})
	eng.World.RegisterSystem(&engine.InputSystem{})
	eng.World.RegisterSystem(&engine.AnimationSystem{})
	eng.World.RegisterSystem(&engine.PhysicsSystem{})
	eng.World.RegisterSystem(&engine.MovementSystem{})
	eng.World.RegisterSystem(&engine.EntityCollisionSystem{})
	eng.World.RegisterSystem(&engine.EntityCollectionSystem{})
	eng.World.RegisterSystem(&engine.AudioSystem{})
	eng.World.RegisterSystem(&engine.HudTextSystem{})

	// order matters here
	eng.World.RegisterSystem(&engine.MapRenderSystem{})
	eng.World.RegisterSystem(&engine.RenderSystem{})
	eng.World.RegisterSystem(&engine.TextRenderSystem{})

	eng.World.RegisterEntityBuilder("coin", CreateCoin)
	eng.World.RegisterEntityBuilder("player", CreatePlayer)
	eng.World.RegisterEntityBuilder("heart", CreateHeart)

	eng.Map.Load("level2", eng.World)
	CreateScoreHud(eng.World)
	CreateHealthHud(eng.World)
	// eng.Audio.PlayBgMusic("pogs.mp3")

	eng.Run()
}
