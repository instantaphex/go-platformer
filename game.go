package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

type Game struct {
	running bool
	renderer *sdl.Renderer
	window *sdl.Window
	keysHeld map[sdl.Keycode]bool
}

func (g *Game) Run() int {
	if g.Init() == false {
		return -1
	}

	g.window.UpdateSurface()
	g.running = true
	for g.running {
		g.HandleInput()
		g.Update()
		g.Render()
	}
	g.Cleanup()
	return 0
}

func (g *Game) Init() bool {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
		return false
	}

	g.window, err = sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WWIDTH,
		WHEIGHT,
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
		return false
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create rederer: %s\n", err)
		return false
	}

	EntityList = append(EntityList, NewEntity("/Users/jb/go/src/github.com/instantaphex/pastor-carrol/yoshi.png", 64, 64, 8))
	areaControl.Load("./maps/1.area")
	g.keysHeld = make(map[sdl.Keycode]bool)
	return true
}

func (g *Game) HandleInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		g.Event(event)
	}
}

func (g *Game) Event(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		println("Quit")
		g.running = false
		break
	case *sdl.KeyboardEvent:
		sym := e.(*sdl.KeyboardEvent).Keysym.Sym
		if t.State == sdl.PRESSED {
			g.keysHeld[sym] = true
		}
		if t.State == sdl.RELEASED {
			g.keysHeld[sym] = false
		}
	}
}

func (g *Game) Update() {
	fpsControl.Update()
	fmt.Println(fpsControl.GetFps())
	for _, entity := range EntityList {
		entity.Update()
	}
	cameraControl.Update(g.keysHeld)
}

func (g *Game) Render() {
	g.renderer.Clear()
	areaControl.Render(cameraControl.GetX(), cameraControl.GetY())
	for _, entity := range EntityList {
		entity.Render()
	}
	g.renderer.Present()
}

func (g *Game) Cleanup() {
	for _, entity := range EntityList {
		entity.Cleanup()
	}
	areaControl.Cleanup()
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
