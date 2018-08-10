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

type GameEntity interface {
	Update()
	Render()
	Cleanup()
	GetEntity() *Entity
}

type EntityCollision struct {
	entityA *Entity
	entityB *Entity
}

type EntityState struct {
	Asset string
	FrameRate int
	FlipHorizontal bool
	FlipVertical bool
	InheritFlip bool
}

type Entity struct {
	animationControl Animation
	Image *Sprite

	X float32
	Y float32

	Flags int32

	MoveLeft bool
	MoveRight bool
	Dead bool

	EntityType int

	SpeedX float32
	SpeedY float32

	AccelX float32
	AccelY float32

	MaxSpeedX float32
	MaxSpeedY float32

	IsJumping bool

	StateMap map[int]EntityState
}

func NewEntity(stateMap map[int]EntityState, x, y int32) *Entity {
	return &Entity{
		X: float32(x),
		Y: float32(y),
		Image: NewSprite(stateMap[ENTITY_STATE_IDLE].Asset, stateMap[ENTITY_STATE_IDLE].FrameRate),
		MoveLeft: false,
		MoveRight: false,
		EntityType: ENTITY_TYPE_NPC,
		Dead: false,
		Flags: ENTITY_FLAG_GRAVITY,
		SpeedX: 0,
		SpeedY: 0,
		AccelX: 0,
		AccelY: 0,
		MaxSpeedX: 3,
		MaxSpeedY: 8,
		StateMap: stateMap,
	}
}

func (e* Entity) Update() {
	e.UpdateState()
	e.UpdatePosition()
}

func (e *Entity) UpdateState() {
	if !e.MoveLeft && !e.MoveRight {
		e.SetState(ENTITY_STATE_IDLE)
		e.StopMove()
	}
	if e.MoveLeft {
		e.SetState(ENTITY_STATE_LEFT)
		e.AccelX = -0.2
	} else if e.MoveRight {
		e.SetState(ENTITY_STATE_RIGHT)
		e.AccelX = 0.2
	}

	if e.IsJumping {
		e.SetState(ENTITY_STATE_JUMP)
		return
	}
}

func (e *Entity) UpdatePosition() {
	if e.Flags & ENTITY_FLAG_GRAVITY != 0 {
		e.AccelY = .55
	}

	e.SpeedX += e.AccelX * fpsControl.GetSpeedFactor()
	e.SpeedY += e.AccelY * fpsControl.GetSpeedFactor()

	if e.SpeedX > e.MaxSpeedX { e.SpeedX = e.MaxSpeedX }
	if e.SpeedX < -e.MaxSpeedX { e.SpeedX = -e.MaxSpeedX }
	if e.SpeedY > e.MaxSpeedY { e.SpeedY = e.MaxSpeedY }
	if e.SpeedY < -e.MaxSpeedY { e.SpeedY = -e.MaxSpeedY }

	e.Move(e.SpeedX, e.SpeedY)
}

func (e *Entity) SetState(state int) {
	newState := e.StateMap[state]
	if !newState.InheritFlip {
		e.Image.FlipHorizontal = newState.FlipHorizontal
		e.Image.FlipVertical = newState.FlipVertical
	}
	e.Image.SetState(newState)
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
		if e.Flags & ENTITY_FLAG_GHOST != 0 {
			e.PosValid(int32(e.X + newX), int32(e.Y + newY))
			e.X += newX
			e.Y += newY
		} else {
			if e.PosValid(int32(e.X + newX), int32(e.Y)) {
				e.X += newX
			} else {
				e.SpeedX = 0
			}

			if e.PosValid(int32(e.X), int32(e.Y + newY)) {
				e.Y += newY
			} else {
				e.SpeedY = 0
				if moveY > 0 {
					e.IsJumping = false
				}
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
	if e.SpeedX > 0 {
		e.AccelX = -.5
	}

	if e.SpeedX < 0 {
		e.AccelX = .5
	}

	if e.SpeedX < .2 && e.SpeedX > -.2 {
		e.AccelX = 0
		e.SpeedX = 0
	}
}

func (e *Entity) Jump() bool {
	e.SetState(ENTITY_STATE_JUMP)
	e.IsJumping = true
	e.SpeedY = -e.MaxSpeedY
	audioManager.PlaySoundEffect("jump.wav")
	return true
}

func (e *Entity) OnCollision(entity *Entity) bool {
	return true
}

func (e *Entity) Collides(oX int32, oY int32, oW int32, oH int32) bool {
	var left1, left2, right1, right2, top1, top2, bottom1, bottom2 int32

	tX := int32(e.X)
	tY := int32(e.Y)

	left1 = tX
	left2 = oX

	right1 = left1 + e.Image.W - 1
	right2 = oX + oW - 1

	top1 = tY
	top2 = oY

	bottom1 = top1 + e.Image.H - 1
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

	if e.Flags & ENTITY_FLAG_MAPONLY != 0 {

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

func (e *Entity) PosValidEntity(gameEntity GameEntity, newX int32, newY int32) bool {
	entity := gameEntity.GetEntity()
	if e != entity &&
	   entity != nil &&
	   !entity.Dead &&
	   entity.Flags & ENTITY_FLAG_MAPONLY == 0 &&
	   entity.Collides(newX, newY, e.Image.W - 1, e.Image.H - 1) {
		entityCollision := EntityCollision{entityA: entity, entityB: e}
		EntityCollisionList = append(EntityCollisionList, entityCollision)
		return false
	}

	return true
}

func (e *Entity) GetEntity() *Entity {
	return e
}

func (e* Entity) Render() {
	x := int32(e.X - cameraControl.GetX())
	y := int32(e.Y - cameraControl.GetY())
	e.Image.Render(x, y)
}

func (e* Entity) Cleanup() {
}
