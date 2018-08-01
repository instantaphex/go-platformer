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

type EntityCollision struct {
	entityA *Entity
	entityB *Entity
}

type Entity struct {
	animationControl Animation
	texture *sdl.Texture

	X float32
	Y float32
	W int32
	H int32

	flags int32

	moveLeft bool
	moveRight bool
	dead bool
	canJump bool

	entType int
	AnimState int

	speedX float32
	speedY float32

	accelX float32
	accelY float32

	maxSpeedX float32
	maxSpeedY float32

	currentFrameCol int32
	currentFrameRow int32

	colX int32
	colY int32
	colWidth int32
	colHeight int32
}

func NewEntity(sheet string, width int32, height int32, maxFrames int) *Entity {
	txt, err := gfx.Load(game.renderer, sheet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load entity sprite sheet: %s\n", err)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get texture size: %s\n", err)
	}
	return &Entity{
		X: 0,
		Y: 0,
		W: width,
		H: height,
		AnimState: 0,
		texture: txt,
		animationControl: NewAnimation(maxFrames),
		moveLeft: false,
		moveRight: false,
		entType: ENTITY_TYPE_NPC,
		dead: false,
		flags: ENTITY_FLAG_GRAVITY,
		speedX: 0,
		speedY: 0,
		accelX: 0,
		accelY: 0,
		maxSpeedX: 10,
		maxSpeedY: 12,
		currentFrameCol: 0,
		currentFrameRow: 0,
		colX: 0,
		colY: 0,
		colWidth: 0,
		colHeight: 0,
	}
}

func (e* Entity) Load(file string, width int, height int, maxFrames int) bool {
	return true
}

func (e* Entity) Update() {
	if !e.moveLeft && !e.moveRight {
		e.StopMove()
	}

	if e.moveLeft {
		e.accelX = -0.5
	} else if e.moveRight {
		e.accelX = 0.5
	}

	if e.flags & ENTITY_FLAG_GRAVITY != 0 {
		e.accelY = .75
	}

	e.speedX += e.accelX * fpsControl.GetSpeedFactor()
	e.speedY += e.accelY * fpsControl.GetSpeedFactor()

	if e.speedX > e.maxSpeedX { e.speedX = e.maxSpeedX }
	if e.speedX < -e.maxSpeedX { e.speedX = -e.maxSpeedX }
	if e.speedY > e.maxSpeedY { e.speedY = e.maxSpeedY }
	if e.speedY < -e.maxSpeedY { e.speedY = -e.maxSpeedY }

	e.Animate()
	e.Move(e.speedX, e.speedY)
}

func (e *Entity) Animate() {
	if e.moveLeft {
		e.currentFrameCol = 0
	} else if e.moveRight {
		e.currentFrameCol = 1
	}
	e.animationControl.Animate()
}

func (e *Entity) Move(moveX float32, moveY float32) {
	if moveX == 0 && moveY == 0 {
		return
	}

	var newX float32 = 0.0
	var newY float32 = 0.0

	moveX *= fpsControl.GetSpeedFactor()
	moveY *= fpsControl.GetSpeedFactor()

	if moveX != 0 {
		if moveX >= 0 {
			newX = fpsControl.GetSpeedFactor()
		} else {
			newX = -fpsControl.GetSpeedFactor()
		}
	}

	if moveY != 0 {
		if moveY >= 0 {
			newY = fpsControl.GetSpeedFactor()
		} else {
			newY = -fpsControl.GetSpeedFactor()
		}
	}

	for {
		if e.flags & ENTITY_FLAG_GHOST != 0 {
			e.PosValid(int32(e.X + newX), int32(e.Y + newY))
			e.X += newX
			e.Y += newY
		} else {
			if e.PosValid(int32(e.X + newX), int32(e.Y)) {
				e.X += newX
			} else {
				e.speedX = 0
			}

			if e.PosValid(int32(e.X), int32(e.Y + newY)) {
				e.Y += newY
			} else {
				if moveY > 0 {
					e.canJump = true
				}
				e.speedY = 0
			}
		}

		moveX += -newX
		moveY += -newY

		if newX > 0 && moveX <= 0 { newX = 0 }
		if newX < 0 && moveX >= 0 { newX = 0 }

		if newY > 0 && moveY <= 0 { newY = 0 }
		if newY < 0 && moveY >= 0 { newY = 0 }

		if moveX == 0 { newX = 0 }
		if moveY == 0 { newY = 0 }

		if moveX == 0 && moveY == 0 { break }
		if newX == 0 && newY == 0 { break }
	}
}

func (e *Entity) StopMove() {
	if e.speedX > 0 {
		e.accelX = -1
	}

	if e.speedX < 0 {
		e.accelX = 1
	}

	if e.speedX < 2.0 && e.speedX > -2.0 {
		e.accelX = 0
		e.speedX = 0
	}
}

func (e *Entity) Jump() bool {
	if !e.canJump { return false }
	e.speedY = -e.maxSpeedY
	return true
}

func (e *Entity) OnCollision(entity *Entity) bool {
	e.Jump()
	return true
}

func (e *Entity) Collides(oX int32, oY int32, oW int32, oH int32) bool {
	var left1, left2, right1, right2, top1, top2, bottom1, bottom2 int32

	tX := int32(e.X) + e.colX
	tY := int32(e.Y) + e.colY

	left1 = tX
	left2 = oX

	right1 = left1 + e.W - 1 - e.colWidth
	right2 = oX + oW - 1

	top1 = tY
	top2 = oY

	bottom1 = top1 + e.H - 1 - e.colHeight
	bottom2 = oY + oH - 1

	if bottom1 < top2 { return false }
	if top1 > bottom2 { return false }
	if right1 < left2 { return false }
	if left1 > right2 { return false }

	return true
}

func (e *Entity) PosValid(newX int32, newY int32) bool {
	retVal := true
	startX := (newX + e.colX) / TILE_SIZE
	startY := (newY + e.colY) / TILE_SIZE

	endX := ((newX + e.colX) + e.W - e.colWidth - 1) / TILE_SIZE
	endY := ((newY + e.colY) + e.H - e.colHeight - 1) / TILE_SIZE

	for iY := startY; iY <= endY; iY++ {
		for iX := startX; iX <= endX; iX++ {
			tile := areaControl.GetTile(iX * TILE_SIZE, iY * TILE_SIZE)

			if e.PosValidTile(tile) == false {
				retVal = false
			}
		}
	}

	if e.flags & ENTITY_FLAG_MAPONLY != 0 {

	} else {
		for i := 0; i < len(EntityList); i++ {
			if e.PosValidEntity(EntityList[i], newX, newY) == false {
				retVal = false
			}
		}
	}
	return retVal
}

func (e *Entity) PosValidTile(tile *Tile) bool {
	if tile == nil {
		return true
	}
	if tile.TypeID == TILE_TYPE_BLOCK {
		return false
	}
	return true
}

func (e *Entity) PosValidEntity(entity *Entity, newX int32, newY int32) bool {
	if e != entity &&
		entity != nil &&
		!entity.dead &&
		entity.flags ^ ENTITY_FLAG_MAPONLY == 0 &&
		entity.Collides(newX + e.colX, newY + e.colY, e.W - e.colWidth - 1, e.H - e.colHeight - 1) {
		entityCollision := EntityCollision{}
		EntityCollisionList = append(EntityCollisionList, entityCollision)
		return false
	}
	return true
}

func (e* Entity) Render() {
gfx.DrawPart(
		game.renderer,
		e.texture,
		int32(e.X - cameraControl.GetX()),
		int32(e.Y - cameraControl.GetY()),
		e.currentFrameCol * e.W,
		(e.currentFrameRow + int32(e.animationControl.GetCurrentFrame())) * e.H,
		e.W,
		e.H,
	)
}

func (e* Entity) Cleanup() {
	e.texture.Destroy()
}
