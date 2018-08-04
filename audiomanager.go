package main

import (
	"github.com/veandco/go-sdl2/mix"
	"fmt"
	"os"
)

type AudioManager struct {
	Sounds map[string]*mix.Music
}

func (mm *AudioManager) Init() {
	if err := mix.Init(mix.INIT_MP3); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing audio: %s\n", err)
	}

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		fmt.Println(err)
		return
	}
	fileNames, err := fileManager.GetDirectoryContents("audio")
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
	path := fileManager.GetAudioPath(asset)
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
	path := fileManager.GetAudioPath(name)
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

