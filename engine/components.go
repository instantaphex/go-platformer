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
	Asset       string
	FrameRate   int
	Flip        sdl.RendererFlip
	Infinite    bool
	Orientation Orientation
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
	X float32
	Y float32
	W int32
	H int32
}

func (t *Transform) GetBB() sdl.Rect {
	return sdl.Rect {
		X: int32(t.X - 1),
		Y: int32(t.Y - 1),
		W: int32(t.W),
		H: int32(t.H),
	}
}

func (t *Transform) GetPotentialBB(x, y int32) sdl.Rect {
	return sdl.Rect {
		X: x,
		Y: y,
		W: int32(t.W),
		H: int32(t.H),
	}
}

type Velocity struct {
	SpeedX    float32
	SpeedY    float32
	AccelX    float32
	AccelY    float32
	MaxSpeedX float32
	MaxSpeedY float32
}

type State struct {
	Jumping     bool
	CanJump     bool
	Grounded    bool
	MoveRight   bool
	MoveLeft    bool
	Rolling     bool
	Shooting    bool
	Orientation Orientation
	Flip        sdl.RendererFlip
	State       StateKey
}


type Animation struct {
	AnimationStates map[StateKey]AnimationState
	AnimState       StateKey
	CurrentFrame    int
	FrameInc        int
	FrameRate       int
	OldTime         uint32
	MaxFrames       int
	Complete        bool
}

func (a *Animation) CurrentState() AnimationState {
	return a.AnimationStates[a.AnimState]
}

// empty components for tagging
type Controller struct {}
type Focused struct {}
type Collidable struct {}

