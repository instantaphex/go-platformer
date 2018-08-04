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
	player *Player
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

	g.renderer.SetScale(SCALE, SCALE)

	audioManager.Init()
	audioManager.PlayBgMusic("theme.mp3")

	/*  DUMPING GROUND */
	textureAtlas.Init()
	g.player = NewPlayer(200, 200)
	EntityList = append(EntityList, g.player)
	mapControl.Load("testmap")
	cameraControl.targetMode = TARGET_MODE_CENTER
	cameraControl.SetTarget(g.player)
	/*  DUMPING GROUND */
	return true
}

func (g *Game) Run() int {
	if g.Init() == false {
		return -1
	}

	g.window.UpdateSurface()
	g.running = true
	for g.running {
		g.HandleEvents()
		g.Update()
		g.Render()
	}
	g.Cleanup()
	return 0
}

func (g *Game) HandleEvents() {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			g.running = false
			break
		case *sdl.KeyboardEvent:
			inputManager.OnKeyboardEvent(e.(*sdl.KeyboardEvent))
			sym := e.(*sdl.KeyboardEvent).Keysym.Sym
			inputManager.KeysHeld[sym] = t.State == sdl.PRESSED
		}
	}
}

func (g *Game) Update() {
	fpsControl.Update()
	for _, entity := range EntityList {
		entity.Update()
	}
	for i := 0; i < len(EntityCollisionList); i++ {
		a := EntityCollisionList[i].entityA
		b := EntityCollisionList[i].entityB

		if a == nil || b == nil { continue }

		if a.OnCollision(b) {
			b.OnCollision(a)
		}
	}
	EntityCollisionList = nil

	// cameraControl.Update(inputManager.KeysHeld)
}

func (g *Game) Render() {
	g.renderer.Clear()
	mapControl.Render(int32(-cameraControl.GetX()), int32(-cameraControl.GetY()))
	for _, entity := range EntityList {
		entity.Render()
	}
	g.renderer.Present()
}

func (g *Game) Cleanup() {
	for _, entity := range EntityList {
		entity.Cleanup()
	}
	mapControl.Cleanup()
	textureAtlas.Cleanup()
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
