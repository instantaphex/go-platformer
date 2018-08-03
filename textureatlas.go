package main

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/veandco/go-sdl2/sdl"
)

type AnimationFrames struct {
	Frames map[string][]AnimationFrame `json:"frames"`
	Texture *sdl.Texture
}

type AnimationFrame struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
	SourceX int32 `json:"sourceX"`
	SourceY int32 `json:"sourceY"`
	SourceW int32 `json:"sourceW"`
	SourceH int32 `json:"sourceH"`
	Rotated bool `json:"rotated"`
}

func (af *AnimationFrames) Parse() *AnimationFrames {
	jsonFile, err := os.Open(fileManager.GetPath("assets", "assets", "json"))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var frames AnimationFrames
	json.Unmarshal(byteValue, &frames)
	af.Texture, err = gfx.Load(game.renderer, fileManager.GetImagePath("assets"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load assets: %s\n", err)
	}
	return &frames
}

