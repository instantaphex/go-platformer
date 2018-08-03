package main

import "github.com/veandco/go-sdl2/sdl"

type Renderable interface {
	Render(x, y int32)
}

type Sprite struct {
	X int32
	Y int32
	W int32
	H int32
	FlipHorizontal bool
	FlipVertical bool
	frames []AnimationFrame
	animation Animation
}

func NewSprite(asset string) *Sprite {
	animationFrames := textureAtlas.Get(asset)
	firstFrame := animationFrames[0]
	return &Sprite {
		X: firstFrame.X,
		Y: firstFrame.Y,
		W: firstFrame.W,
		H: firstFrame.H,
		FlipHorizontal: false,
		FlipVertical: false,
		frames: textureAtlas.Get(asset),
		animation: NewAnimation(len(animationFrames)),
	}
}

func (s *Sprite) Render(x, y int32) {
	frame := s.frames[s.animation.GetCurrentFrame()]
	xOffset := frame.SourceW - frame.W
	yOffset := frame.SourceH - frame.H
	s.X = x
	s.Y = y
	s.W = frame.W + xOffset
	s.H = frame.H + yOffset
	renderFlip := sdl.FLIP_NONE
	if s.FlipHorizontal {
		renderFlip |= sdl.FLIP_HORIZONTAL
	}
	if s.FlipVertical {
		renderFlip |= sdl.FLIP_VERTICAL
	}
	gfx.DrawPart(game.renderer, textureAtlas.Texture, x + xOffset, y + yOffset, frame.X, frame.Y, frame.W, frame.H, renderFlip)
	s.animation.Animate()
}

func (s *Sprite) SetAsset(asset string) {
	s.frames = textureAtlas.Get(asset)
}

