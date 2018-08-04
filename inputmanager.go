package main

import "github.com/veandco/go-sdl2/sdl"

type InputManager struct {
	KeysHeld map[sdl.Keycode]bool
}

func NewInputManager() *InputManager {
	return &InputManager {
		KeysHeld: make(map[sdl.Keycode]bool),
	}
}

func (i *InputManager) OnKeyboardEvent(evt *sdl.KeyboardEvent) {
	t := evt
	sym := evt.Keysym.Sym
	if t.State == sdl.PRESSED {
		inputManager.KeysHeld[sym] = true
	}
	if t.State == sdl.RELEASED {
		inputManager.KeysHeld[sym] = false
	}
}
