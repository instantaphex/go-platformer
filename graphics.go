package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"os"
)

type Graphics struct {
}

func (g *Graphics) Load(r *sdl.Renderer, file string) (*sdl.Texture, error) {
	i, err := img.LoadTexture(r, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", err)
		return nil, err
	}
	return i, err
}

func (g *Graphics) Draw(r *sdl.Renderer, texture *sdl.Texture, x int32, y int32) {
	_, _, w, h, err := texture.Query()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to draw texture: %s\n", err)
		panic(err)
	}
	g.DrawPart(r, texture, x, y, 0, 0, w, h)
}

func (g *Graphics) DrawPart(r *sdl.Renderer, texture *sdl.Texture, x int32, y int32, clipX int32, clipY int32, w int32, h int32, flip ...sdl.RendererFlip) {
	src := sdl.Rect{clipX, clipY, w, h}
	dst := sdl.Rect{x, y, w, h}
	// r.Copy(texture, &src, &dst)
	renderFlip := sdl.FLIP_NONE
	if (len(flip) > 0) {
		renderFlip = flip[0]
	}
	r.CopyEx(texture, &src, &dst, 0.0, nil, renderFlip)
}
