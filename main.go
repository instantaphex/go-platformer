package main

import "github.com/instantaphex/platformer/engine"

var game = Game{}
var gfx = Graphics{}
var EntityList []GameEntity
var EntityCollisionList []EntityCollision
var cameraControl = Camera{}
var mapControl = Map{}
var fpsControl = Fps{}
var fileManager = FileManager{}
var textureAtlas = TextureAtlas{}
var audioManager = AudioManager{}
var inputManager = NewInputManager()

func main() {
	// game.Run()
	eng := engine.New(engine.EngineConfig{
		WindowWidth: 1024,
		WindowHeight: 768,
		WindowTitle: "Platformer",
		Scale: 2,
		DrawDebug: false,
	})
	eng.Run()
}
