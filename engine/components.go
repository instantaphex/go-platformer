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
	COMPONENT_TAG = 1 << 9
	COMPONENT_INVENTORY = 1 << 11
	COMPONENT_COLLECTIBLE = 1 << 12
	COMPONENT_TEXT = 1 << 13
	COMPONENT_HUD = 1 << 14
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
	ENTITY_STATE_WALLR
	ENTITY_STATE_WALLL
)

type Transform struct {
	X float32
	Y float32
	W int32
	H int32

	SpeedX    float32
	SpeedY    float32
	AccelX    float32
	AccelY    float32
	MaxSpeedX float32
	MaxSpeedY float32

	Sensor struct {
		Top bool
		Bottom bool
		Left bool
		Right bool
	}
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

func (t *Transform) GetSensorPoints() (sdl.Point, sdl.Point, sdl.Point, sdl.Point) {
	top := sdl.Point{ X: int32(t.X) + (t.W / 2), Y: int32(t.Y) - 1}
	bottom := sdl.Point{ X: int32(t.X) + (t.W / 2), Y: (int32(t.Y) + t.H) + 1 }
	left := sdl.Point{ X: int32(t.X) - 1, Y: int32(t.Y) + (t.H / 2) }
	right := sdl.Point{ X: (int32(t.X) + t.W) + 1, Y: int32(t.Y) + (t.H / 2)}
	return top, bottom, left, right
}

type State struct {
	Jumping     bool
	CanJump     bool
	Grounded    bool
	MoveRight   bool
	MoveLeft    bool
	Rolling     bool
	Shooting    bool
	LeftSlide   bool
	RightSlide  bool
	Sliding     bool
	JumpCount   int
	JumpFrameCount int
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

type Tag struct {
	Value string
}

type Inventory struct {
	Items map[string]int
}

type Collectible struct {
	Type string
	Value int
}

type Text struct {
	Value string
}

// empty components for tagging
type Controller struct {}
type Focused struct {}
type Collidable struct {}
type Hud struct {}

