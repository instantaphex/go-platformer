package main

import (
	"os"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"encoding/json"
)

type TextureAtlas struct {
	Texture *sdl.Texture
	Images map[string][]AnimationFrame `json:"frames"`
}

type AnimationFrames struct {
	Images map[string][]AnimationFrame `json:"frames"`
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

func (ta *TextureAtlas) Init() {
	ta.LoadAssetJson()
	ta.LoadTexture()
}

func (ta* TextureAtlas) LoadAssetJson() {
	jsonFile, err := os.Open(fileManager.GetPath("assets", "assets", "json"))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, ta)
}

func (ta *TextureAtlas) LoadTexture() {
	var err error
	ta.Texture, err = gfx.Load(game.renderer, fileManager.GetImagePath("assets"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load assets: %s\n", err)
	}
}

func (ta *TextureAtlas) Get(asset string) []AnimationFrame {
	return ta.Images[asset]
}
