package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	KeyStateUp = iota
	KeyStateDown
	KeyStateJustDown
	KeyStateJustUp
)

type KeyState struct {
	lastState bool
	currentState bool
}

func (key *KeyState) Set(state bool) {
	key.lastState = key.currentState
	key.currentState = state
}

func (key *KeyState) State() int {
	if key.lastState {
		if key.currentState {
			return KeyStateDown
		}
		return KeyStateJustUp
	}
	if key.currentState {
		return KeyStateJustDown
	}
	return KeyStateUp
}

// key was just pressed
func (key *KeyState) JustPressed() bool {
	return !key.lastState && key.currentState
}

// key was just released
func (key *KeyState) JustReleased() bool {
	return key.lastState && !key.currentState
}

// not being pressed
func (key *KeyState) Up() bool {
	return !key.lastState && !key.currentState
}

// currently held
func (key *KeyState) Down() bool {
	return key.lastState && key.currentState
}

type InputManager struct {
	KeysHeld map[sdl.Keycode]bool
	KeyStates map[sdl.Keycode]*KeyState
	PressedListeners map[sdl.Keycode][]func()
	ReleasedListeners map[sdl.Keycode][]func()
}

func NewInputManager() *InputManager {
	return &InputManager {
		KeysHeld: make(map[sdl.Keycode]bool),
		KeyStates: make(map[sdl.Keycode]*KeyState),
		PressedListeners: make(map[sdl.Keycode][]func()),
		ReleasedListeners: make(map[sdl.Keycode][]func()),
	}
}

func (i *InputManager) UpdateKeyStates() {
	for _, state := range i.KeyStates {
		state.Set(state.currentState)
	}
}

func (i *InputManager) OnKeyboardEvent(evt *sdl.KeyboardEvent) {
	t := evt
	sym := evt.Keysym.Sym

	if t.State == sdl.PRESSED {
		i.KeysHeld[sym] = true
		i.SetKeyState(sym, true)
		// i.executeCallbacks(&i.PressedListeners, sym)
	}

	if t.State == sdl.RELEASED {
		i.KeysHeld[sym] = false
		i.SetKeyState(sym, false)
		// i.executeCallbacks(&i.ReleasedListeners, sym)
	}
}

func (i *InputManager) KeyState(key sdl.Keycode) *KeyState {
	if _, ok := i.KeyStates[key]; !ok {
		i.KeyStates[key] = &KeyState{}
		i.KeyStates[key].Set(false)
	}
	return i.KeyStates[key]
}

func (i *InputManager) SetKeyState(key sdl.Keycode, state bool) {
	if _, ok := i.KeyStates[key]; !ok {
		i.KeyStates[key] = &KeyState{}
	}
	i.KeyStates[key].Set(state)
}

/*func (i *InputManager) RegisterKeyListener(sym sdl.Keycode, evtType string, fn func()) Subscription {
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
}*/

