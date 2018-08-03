package main

const (
	ENTITY_TYPE_NPC = iota
	ENTITY_TYPE_PLAYER
)

const (
	ENTITY_FLAG_NONE = 0
	ENTITY_FLAG_GRAVITY = 0x00000001
	ENTITY_FLAG_GHOST = 0x00000002
	ENTITY_FLAG_MAPONLY = 0x00000004
)

const (
	ENTITY_STATE_IDLE = iota
	ENTITY_STATE_LEFT
	ENTITY_STATE_RIGHT
	ENTITY_STATE_JUMP
	ENTITY_STATE_SHOOT
	ENTITY_STATE_DIE
)

type EntityCollision struct {
	entityA *Entity
	entityB *Entity
}

type Entity struct {
	animationControl Animation
	Image *Sprite

	X float32
	Y float32

	flags int32

	moveLeft bool
	moveRight bool
	dead bool
	canJump bool

	entType int

	speedX float32
	speedY float32

	accelX float32
	accelY float32

	maxSpeedX float32
	maxSpeedY float32

	StateChannel chan int
	CurrentState int
}

func NewEntity(asset string) *Entity {
	return &Entity{
		X: 0,
		Y: 0,
		Image: NewSprite(asset),
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
		StateChannel: make(chan int),
	}
}

func (e* Entity) Load(file string, width int, height int, maxFrames int) bool {
	return true
}

func (e* Entity) Update() {
	if !e.moveLeft && !e.moveRight {
		e.SetState(ENTITY_STATE_IDLE)
		e.StopMove()
	}
	if e.moveLeft {
		e.SetState(ENTITY_STATE_LEFT)
		e.accelX = -0.5
	} else if e.moveRight {
		e.SetState(ENTITY_STATE_RIGHT)
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

	// e.Animate()
	e.Move(e.speedX, e.speedY)
}

func (e *Entity) SetState(state int) {
	var newState string
	if state == ENTITY_STATE_RIGHT {
		newState = "Player/Run"
		e.Image.FlipHorizontal = false
	}
	if state == ENTITY_STATE_LEFT {
		newState = "Player/Run"
		e.Image.FlipHorizontal = true
	}
	if state == ENTITY_STATE_IDLE {
		newState = "Player/Idle"
	}
	if state == ENTITY_STATE_JUMP {
		newState = "Player/Fall-Jump-WallJ/Jump"
	}
	e.Image.SetAsset(newState)
	e.CurrentState = state
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
	e.SetState(ENTITY_STATE_JUMP)
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

	tX := int32(e.X)
	tY := int32(e.Y)

	left1 = tX
	left2 = oX

	right1 = left1 + e.Image.W - 1//  - e.colWidth
	right2 = oX + oW - 1

	top1 = tY
	top2 = oY

	bottom1 = top1 + e.Image.H - 1//  - e.colHeight
	bottom2 = oY + oH - 1

	if bottom1 < top2 { return false }
	if top1 > bottom2 { return false }
	if right1 < left2 { return false }
	if left1 > right2 { return false }

	return true
}

func (e *Entity) PosValid(newX int32, newY int32) bool {
	retVal := true
	startX := (newX) / TILE_SIZE
	startY := (newY) / TILE_SIZE

	endX := ((newX) + e.Image.W - 1) / TILE_SIZE
	endY := ((newY) + e.Image.H - 1) / TILE_SIZE

	for iY := startY; iY <= endY; iY++ {
		for iX := startX; iX <= endX; iX++ {
			tile := mapControl.GetTile(iX * TILE_SIZE, iY * TILE_SIZE)

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
		entity.Collides(newX, newY, e.Image.W - 1, e.Image.H - 1) {
		entityCollision := EntityCollision{}
		EntityCollisionList = append(EntityCollisionList, entityCollision)
		return false
	}
	return true
}

func (e* Entity) Render() {
	x := int32(e.X - cameraControl.GetX())
	y := int32(e.Y - cameraControl.GetY())
	e.Image.Render(x, y)
}

func (e* Entity) Cleanup() {
}
