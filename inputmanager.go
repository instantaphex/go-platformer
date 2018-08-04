package main

import "github.com/veandco/go-sdl2/sdl"

type InputManager struct {
	KeysHeld map[sdl.Keycode]bool
	PressedListeners map[sdl.Keycode][]func()
	ReleasedListeners map[sdl.Keycode][]func()
}

type Subscription struct {
	Unsubscribe func()
}

func NewInputManager() *InputManager {
	return &InputManager {
		KeysHeld: make(map[sdl.Keycode]bool),
		PressedListeners: make(map[sdl.Keycode][]func()),
		ReleasedListeners: make(map[sdl.Keycode][]func()),
	}
}

func (i *InputManager) OnKeyboardEvent(evt *sdl.KeyboardEvent) {
	t := evt
	sym := evt.Keysym.Sym

	if t.State == sdl.PRESSED {
		inputManager.KeysHeld[sym] = true
		i.executeCallbacks(&i.PressedListeners, sym)
	}

	if t.State == sdl.RELEASED {
		inputManager.KeysHeld[sym] = false
		i.executeCallbacks(&i.ReleasedListeners, sym)
	}
}

func (i *InputManager) RegisterKeyListener(sym sdl.Keycode, evtType string, fn func()) Subscription {
	var idx int
	if evtType == "pressed" {
		i.PressedListeners[sym] = append(i.PressedListeners[sym], fn)
		idx = len(i.PressedListeners) - 1
	} else if evtType == "released" {
		i.ReleasedListeners[sym] = append(i.ReleasedListeners[sym], fn)
		idx = len(i.ReleasedListeners) - 1
	}
	return Subscription {
		Unsubscribe: func() {
			if evtType == "pressed" {
				i.PressedListeners[sym] = append(
					i.PressedListeners[sym][:idx],
					i.PressedListeners[sym][idx+1:]...
				)
			}
			if evtType == "released" {
				i.ReleasedListeners[sym] = append(
					i.ReleasedListeners[sym][:idx],
					i.ReleasedListeners[sym][idx+1:]...
				)
			}
		},
	}
}

func (i *InputManager) executeCallbacks(cbMap *map[sdl.Keycode][]func(), key sdl.Keycode) {
	if (*cbMap)[key] != nil {
		for _, fn := range (*cbMap)[key] {
			fn()
		}
	}
}
