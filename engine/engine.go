package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

type Engine struct {
	running bool
	World *World
	renderer *sdl.Renderer
	window *sdl.Window

	Audio *AudioManager
	File *FileManager
	Graphics *Graphics
	Input *InputManager
	Assets *TextureAtlas
	FPS *Fps
	Map *Map
	Camera *Camera
	Config EngineConfig
}

type EngineConfig struct {
	WindowWidth int32
	WindowHeight int32
	WindowTitle string
	Scale float32
	DrawDebug bool
}

func New(cfg EngineConfig) *Engine {
	eng := &Engine{}

	eng.FPS = 		&Fps{}
	eng.File = 		&FileManager{}
	eng.Input = 	NewInputManager()
	eng.Assets = 	&TextureAtlas{engine: eng}
	eng.Graphics = 	&Graphics{engine: eng}
	eng.Map = 		&Map{engine: eng}
	eng.Audio = 	&AudioManager{engine: eng}
	eng.Camera = 	&Camera{engine: eng}
	eng.Camera.targetMode = TARGET_MODE_CENTER
	eng.Config =	cfg

	eng.Init()
	return eng
}

func (g *Engine) Init() error {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
		return err
	}

	g.window, err = sdl.CreateWindow(
		g.Config.WindowTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		g.Config.WindowWidth,
		g.Config.WindowHeight,
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
		return err
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create rederer: %s\n", err)
		return err
	}

	g.renderer.SetScale(g.Config.Scale, g.Config.Scale)

	g.Assets.Init()
	g.Audio.Init()

	/* DUMPING GROUND */
	g.World = &World{}
	g.World.RegisterSystem(&CameraSystem{})
	g.World.RegisterSystem(&InputSystem{})
	g.World.RegisterSystem(&AnimationSystem{})
	g.World.RegisterSystem(&VelocitySystem{})
	g.World.RegisterSystem(&MovementSystem{})

	// order matters here
	g.World.RegisterSystem(&MapRenderSystem{})
	g.World.RegisterSystem(&RenderSystem{})

	// g.World.CreatePlayer(g, 100, 0)
	g.Map.Load("level1", g.World)
	/* DUMPING GROUND */

	return nil
}

func (g *Engine) Run() int {
	g.running = true
	for g.running {
		g.HandleEvents()
		g.Update()
	}
	g.Cleanup()
	return 0
}

func (g *Engine) HandleEvents() {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			g.running = false
			break
		case *sdl.KeyboardEvent:
			g.Input.OnKeyboardEvent(e.(*sdl.KeyboardEvent))
			sym := e.(*sdl.KeyboardEvent).Keysym.Sym
			g.Input.KeysHeld[sym] = t.State == sdl.PRESSED
			if sym == sdl.K_q {
				g.running = false
			}
		}
	}
}

func (g *Engine) Update() {
	g.FPS.Update()
	g.renderer.Clear()
	g.World.Update(g)
	g.renderer.Present()
}

func (g *Engine) Cleanup() {
	g.Assets.Cleanup()
	g.Audio.Cleanup()

	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

