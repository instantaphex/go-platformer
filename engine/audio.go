package engine

import (
	"fmt"
	"github.com/veandco/go-sdl2/mix"
	"os"
)

type AudioManager struct {
	engine *Engine
	Sounds map[string]*mix.Music
	buffer int
}

func (mm *AudioManager) Init() {
	// mm.buffer = 4096
	mm.buffer = 2048
	if err := mix.Init(mix.INIT_MP3); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing audio: %s\n", err)
	}

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, mm.buffer); err != nil {
		fmt.Println(err)
		return
	}
	fileNames, err := mm.engine.File.GetDirectoryContents("audio")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get directory list: %s", err)
	}

	mm.Sounds = make(map[string]*mix.Music)

	for _, file := range fileNames {
		music, err := mix.LoadMUS(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load audio file %s: %s\n", file, err)
		}
		mm.Sounds[file] = music
	}
}

func (mm *AudioManager) PlayBgMusic(asset string) {
	path := mm.engine.File.GetAudioPath(asset)
	music := mm.Sounds[path]
	if music != nil {
		err := music.Play(1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error playing music: %s\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "No sound file named: %s", asset)
	}
}

func (am *AudioManager) PlaySoundEffect(name string) {
	path := am.engine.File.GetAudioPath(name)
	if sound, err := mix.LoadWAV(path); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load sound effect: %s\n", err)
	} else if sound.Play(-1, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to play sound effect: %s", err)
	}
}

func (mm *AudioManager) Cleanup() {
	for _, val := range mm.Sounds {
		val.Free()
	}
	mix.CloseAudio()
	mix.Quit()
}


