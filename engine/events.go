package engine

type Event interface {
	Type() string
	Async() bool
}

type CollisionEvent struct {
	A int
	B int
}
func (ce *CollisionEvent) Type() string { return "collision" }
func (ce *CollisionEvent) Async() bool { return true }

type AudioEvent struct {
	Clip string
}
func (ae *AudioEvent) Type() string { return "audio" }
func (ae *AudioEvent) Async() bool { return true }

type CollectionEvent struct {
	Collectible string
	Total int
	NumCollected int
}
func (ae *CollectionEvent) Type() string { return "collection" }
func (ae *CollectionEvent) Async() bool { return true }

type PhysicsPulseEvent struct {
	SpeedX float32
	SpeedY float32
	AccelX float32
	AccelY float32
	Entity int
}
func (ae *PhysicsPulseEvent) Type() string { return "physics-pulse" }
func (ae *PhysicsPulseEvent) Async() bool { return true }

type PhysicsSetEvent struct {
	SpeedX float32
	SpeedY float32
	AccelX float32
	AccelY float32
	Entity int
}
func (ae *PhysicsSetEvent) Type() string { return "physics-set" }
func (ae *PhysicsSetEvent) Async() bool { return true }
