package main

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
	game.Run()
}
