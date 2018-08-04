package main

import "github.com/veandco/go-sdl2/sdl"

type Animation struct {
	currentFrame int
	frameInc int
	frameRate int
	oldTime uint32

	maxFrames int
}

func NewAnimation(frameCount int, frameRate int) Animation {
	return Animation {
		currentFrame: 0,
		frameInc: 1,
		frameRate: frameRate,
		oldTime: 0,
		maxFrames: frameCount,
	}
}

func (a *Animation) Animate() {
	if a.oldTime + uint32(a.frameRate) > sdl.GetTicks() {
		return
	}

	a.oldTime = sdl.GetTicks()

	a.currentFrame += a.frameInc

	if a.currentFrame >= a.maxFrames {
		a.currentFrame = 0
	}
}

func (a *Animation) SetFrameRate(rate int) {
	a.frameRate = rate
}

func (a *Animation) SetCurrentFrame(frame int) {
	if frame < 0 || frame >= a.maxFrames {
		return
	}

	a.currentFrame = frame
}

func (a *Animation) SetMaxFrames(maxFrames int) {
	a.maxFrames = maxFrames
}

func (a *Animation) GetCurrentFrame() int {
	return a.currentFrame
}
