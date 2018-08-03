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
	player *Entity
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

	g.renderer.SetScale(1, 1)

	/*  DUMPING GROUND */
	textureAtlas.Init()
	g.player = NewPlayer()
	EntityList = append(EntityList, g.player)
	mapControl.Load("testmap")
	g.keysHeld = make(map[sdl.Keycode]bool)
	cameraControl.targetMode = TARGET_MODE_CENTER
	cameraControl.SetTarget(&g.player.X, &g.player.Y)
	/*  DUMPING GROUND */
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
		g.running = false
		break
	case *sdl.KeyboardEvent:
		sym := e.(*sdl.KeyboardEvent).Keysym.Sym
		if t.State == sdl.PRESSED {
			g.keysHeld[sym] = true
			if sym == sdl.K_LEFT {
				g.player.moveLeft = true
			}
			if sym == sdl.K_RIGHT {
				g.player.moveRight = true
			}
			if sym == sdl.K_SPACE {
				g.player.Jump()
			}
		}
		if t.State == sdl.RELEASED {
			g.keysHeld[sym] = false
			if sym == sdl.K_LEFT {
				g.player.moveLeft = false
			}
			if sym == sdl.K_RIGHT {
				g.player.moveRight = false
			}
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

	cameraControl.Update(g.keysHeld)
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
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
