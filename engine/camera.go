package engine

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

	engine *Engine
}

func (c *Camera) Move(moveX, moveY float32) {
	c.x += moveX
	c.y += moveY
}

func (c *Camera) X() float32 {
	if c.targetX != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			ww := float32(c.engine.Config.WindowWidth)
			return *c.targetX - (ww / (2 * (c.engine.Config.Scale)))
		}
		return *c.targetX
	}
	return c.x
}

func (c *Camera) Y() float32 {
	if c.targetY != nil {
		if c.targetMode == TARGET_MODE_CENTER {
			wh := float32(c.engine.Config.WindowHeight)
			return *c.targetY - (wh / (2 * c.engine.Config.Scale))
		}
		return *c.targetY
	}
	return c.y
}

func (c *Camera) SetPos(x, y float32) {
	c.x = x
	c.y = y
}

func (c *Camera) SetTarget(x, y *float32) {
	c.targetX = x
	c.targetY = y
}

func (c *Camera) Update(keysHeld map[sdl.Keycode]bool) {
}
