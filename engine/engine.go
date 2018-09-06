package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

type Engine struct {
	running bool
	paused bool
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
	Events *Dispatcher
	Text *TextManager
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
	eng.Text =      &TextManager{engine: eng}
	eng.Graphics = 	&Graphics{engine: eng}
	eng.Map = 		&Map{engine: eng}
	eng.Audio = 	&AudioManager{engine: eng}
	eng.Camera = 	&Camera{engine: eng}
	eng.Camera.targetMode = TARGET_MODE_CENTER
	eng.Config =	cfg

	eng.Init()
	return eng
}

func (e *Engine) Init() error {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
		return err
	}

	e.window, err = sdl.CreateWindow(
		e.Config.WindowTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		e.Config.WindowWidth,
		e.Config.WindowHeight,
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
		return err
	}

	e.renderer, err = sdl.CreateRenderer(e.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create rederer: %s\n", err)
		return err
	}

	e.renderer.SetScale(e.Config.Scale, e.Config.Scale)

	e.Assets.Init()
	e.Audio.Init()
	e.Text.Init()

	return nil
}

func (g *Engine) Run() int {
	g.running = true
	g.paused = false
	for g.running {
			g.HandleEvents()
		if !g.paused {
			g.Update()
		}
	}
	g.Cleanup()
	return 0
}

func (g *Engine) HandleEvents() {
	/**
	 * set keystates every frame so that lastState and currentState will be set correctly
	 */
	g.Input.UpdateKeyStates()
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
			if sym == sdl.K_p  && t.State == sdl.RELEASED {
				fmt.Println("p pressed")
				g.paused = !g.paused
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
	g.Text.Cleanup()

	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

