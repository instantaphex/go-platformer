package main

var game = Game{}
var gfx = Graphics{}
var EntityList []*Entity
var cameraControl = Camera{}
var areaControl = Area{}
var fpsControl = Fps{}

func main() {
	game.Run()
}
