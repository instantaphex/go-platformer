package main

import "github.com/veandco/go-sdl2/sdl"

const (
	TARGET_MODE_NORMAL = 0 << iota
	TARGET_MODE_CENTER
)

type Camera struct {
	x int32
	y int32

	targetX *int32
	targetY *int32

	targetMode int
}

func (c *Camera) Move(moveX int32, moveY int32) {
	c.x += moveX
	c.y += moveY
}

func (c *Camera) GetX() int32 {
	if c.targetX != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			return *c.targetX - (WWIDTH / 2)
		}
		return *c.targetX
	}
	return c.x
}

func (c *Camera) GetY() int32 {
	if c.targetY != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			return *c.targetY - (WHEIGHT / 2)
		}
		return *c.targetY
	}
	return c.y
}

func (c *Camera) SetPos(x int32, y int32) {
	c.x = x
	c.y = y
}

func (c *Camera) SetTarget(x *int32, y *int32) {
	c.targetX = x
	c.targetY = y
}

func (c *Camera) Update(keysHeld map[sdl.Keycode]bool) {
	if keysHeld[sdl.K_UP] {
		cameraControl.Move(0, 1)
	}
	if keysHeld[sdl.K_DOWN] {
		cameraControl.Move(0, -1)
	}
	if keysHeld[sdl.K_LEFT] {
		cameraControl.Move(1, 0)
	}
	if keysHeld[sdl.K_RIGHT] {
		cameraControl.Move(-1, 0)
	}
}
