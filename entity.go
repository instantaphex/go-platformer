package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

const (
	ENTITY_TYPE_NPC = 0 << iota
	ENTITY_TYPE_PLAYER
)

const (
	ENTITY_FLAG_NONE = 0
	ENTITY_FLAG_GRAVITY = 0x00000001
	ENTITY_FLAG_GHOST = 0x00000002
	ENTITY_FLAG_MAPONLY = 0x00000004
)

type Entity struct {
	animationControl Animation
	texture *sdl.Texture

	X float64
	Y float64
	W int32
	H int32
	AnimState int
}

func NewEntity(sheet string, width int32, height int32, maxFrames int) *Entity {
	txt, err := gfx.Load(game.renderer, sheet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load entity sprite sheet: %s\n", err)
	}
	return &Entity{
		X: 0,
		Y: 0,
		W: width,
		H: height,
		AnimState: 0,
		texture: txt,
		animationControl: NewAnimation(maxFrames),
	}
}

func (e* Entity) Load(file string, width int, height int, maxFrames int) bool {
	return true
}

func (e* Entity) Update() {
	e.animationControl.Animate()
}

func (e* Entity) Render() {
	gfx.DrawPart(game.renderer, e.texture, 290, 220, 0, int32(e.animationControl.GetCurrentFrame()) * e.H, e.W, e.H)
}

func (e* Entity) Cleanup() {
	e.texture.Destroy()
}
