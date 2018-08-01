package main

var game = Game{}
var gfx = Graphics{}
var EntityList []*Entity
var EntityCollisionList []EntityCollision
var cameraControl = Camera{}
var mapControl = Map{}
var fpsControl = Fps{}
var fileManager = FileManager{}

func main() {
	game.Run()
}
