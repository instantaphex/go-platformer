package main

var game = Game{}
var gfx = Graphics{}
var EntityList []*Entity
var EntityCollisionList []EntityCollision
var cameraControl = Camera{}
var areaControl = Area{}
var fpsControl = Fps{}
var fileManager = FileManager{}

func main() {
	game.Run()
}
