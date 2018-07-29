package main

import "github.com/veandco/go-sdl2/sdl"

type Fps struct {
	oldTime uint32
	lastTime uint32
	speedFactor float32
	fps uint32
	frames uint32
}

func (f *Fps) Update() {
	if f.oldTime + 1000 < sdl.GetTicks() {
		f.oldTime = sdl.GetTicks()
		f.fps = f.frames
		f.frames = 0
	}

	speedMilli := sdl.GetTicks() - f.lastTime
	speedSeconds := float32(speedMilli / 1000.0)
	f.speedFactor = speedSeconds * 60.0
	f.lastTime = sdl.GetTicks()
	f.frames++
}

func (f *Fps) GetFps() uint32 {
	return f.fps
}

func (f *Fps) GetSpeedFactor() float32 {
	return f.speedFactor
}
