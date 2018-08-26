package engine

import "github.com/veandco/go-sdl2/sdl"

const (
	COMPONENT_NONE = 0
	COMPONENT_POSITION = 1 << 0
	COMPONENT_VELOCITY = 1 << 1
	COMPONENT_APPEARANCE = 1 << 2
	COMPONENT_ANIMATION = 1 << 3
	COMPONENT_FOCUSED = 1 << 5
	COMPONENT_STATE = 1 << 6
	COMPONENT_COLLIDER = 1 << 6
)

const (
	ORIENTATION_RIGHT = 0
	ORIENTATION_LEFT = 1
)

type Orientation int

type AnimationStateKey int

type AnimationState struct {
	asset string
	frameRate int
	flip sdl.RendererFlip
	infinite bool
	orientation Orientation
}

const (
	ENTITY_STATE_IDLE = iota
	ENTITY_STATE_LEFT
	ENTITY_STATE_RIGHT
	ENTITY_STATE_JUMP
	ENTITY_STATE_SHOOT
	ENTITY_STATE_DIE
	ENTITY_STATE_ROLL
)

type Position struct {
	x float32
	y float32
}

type Collider struct {
	x float32
	y float32
	w int32
	h int32
}

type Velocity struct {
	speedX float32
	speedY float32
	accelX float32
	accelY float32
	maxSpeedX float32
	maxSpeedY float32
}

type Appearance struct {
	flip sdl.RendererFlip
	name string
	frame AnimationFrame
	xOffset int32
	yOffset int32
	w int32
	h int32
}

type State struct {
	animationStates map[AnimationStateKey]AnimationState
	currentAnimKey AnimationStateKey
	desiredAnimKey AnimationStateKey
	animationState AnimationState
	jumping bool
	canJump bool
	grounded bool
	newState bool
	moveRight bool
	moveLeft bool
	rolling bool
	shooting bool
	orientation Orientation
}

type Animation struct {
	currentFrame int
	frameInc int
	frameRate int
	oldTime uint32
	maxFrames int
	complete bool
}

type Focused struct {}

