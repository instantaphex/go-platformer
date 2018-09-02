package engine

type Scene interface {
	Init()
	HandleEvents()
	Update()
	Cleanup()
}
