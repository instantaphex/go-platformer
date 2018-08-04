package main

import "github.com/veandco/go-sdl2/sdl"

const (
	TARGET_MODE_NORMAL = iota
	TARGET_MODE_CENTER
)

type Camera struct {
	x float32
	y float32

	targetX *float32
	targetY *float32

	targetMode int
}

func (c *Camera) Move(moveX float32, moveY float32) {
	c.x += moveX
	c.y += moveY
}

func (c *Camera) GetX() float32 {
	if c.targetX != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			return *c.targetX - (WWIDTH / (2 * SCALE))
		}
		return *c.targetX
	}
	return c.x
}

func (c *Camera) GetY() float32 {
	if c.targetY != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			return *c.targetY - ((WHEIGHT / (2 * SCALE)) + 150)
		}
		return *c.targetY
	}
	return c.y
}

func (c *Camera) SetPos(x float32, y float32) {
	c.x = x
	c.y = y
}

func (c *Camera) SetTarget(entity GameEntity) {
	c.targetX = &entity.GetEntity().X
	c.targetY = &entity.GetEntity().Y
}

func (c *Camera) Update(keysHeld map[sdl.Keycode]bool) {
}
