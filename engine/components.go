package engine

import "github.com/veandco/go-sdl2/sdl"

const (
	COMPONENT_NONE = 0
	COMPONENT_POSITION = 1 << 0
	COMPONENT_VELOCITY = 1 << 1
	COMPONENT_APPEARANCE = 1 << 2
	COMPONENT_ANIMATION = 1 << 3
	COMPONENT_FOCUSED = 1 << 4
	COMPONENT_STATE = 1 << 5
	COMPONENT_CONTROLLER = 1 << 6
	COMPONENT_TRANSFORM = 1 << 7
	COMPONENT_COLLIDABLE = 1 << 8
)

const (
	ORIENTATION_RIGHT = 0
	ORIENTATION_LEFT = 1
)

type Orientation int

type StateKey int

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

type Transform struct {
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

type State struct {
	jumping         bool
	canJump         bool
	grounded        bool
	moveRight       bool
	moveLeft        bool
	rolling         bool
	shooting        bool
	orientation     Orientation
	flip sdl.RendererFlip
	state StateKey
}


type Animation struct {
	animationStates map[StateKey]AnimationState
	animState StateKey
	currentFrame    int
	frameInc        int
	frameRate       int
	oldTime         uint32
	maxFrames       int
	complete        bool
}

// empty components for tagging
type Controller struct {}
type Focused struct {}
type Collidable struct {}

