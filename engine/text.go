package engine

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
)

type TextNode struct {
	value string
	texture *sdl.Texture
}

type TextManager struct {
	engine *Engine
	font *ttf.Font
	cache map[string]*TextNode
}

func (tm *TextManager) Init() {
	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init ttf: %s\n", err)
	}

	fontPath := tm.engine.File.GetFontPath("kongtext")

	var err error
	if tm.font, err = ttf.OpenFont(fontPath, 8); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
	}

	tm.cache = make(map[string]*TextNode)
}

func (tm *TextManager) Write(name string, text string) {
	var err error
	if node, ok := tm.cache[name]; ok {
		if node.value != text {
			node.texture.Destroy()
			node.value = text
			node.texture, err = tm.createText(text)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	} else {
		texture, _ := tm.createText(text)

		tm.cache[name] = &TextNode {
			value: text,
			texture: texture,
		}
	}
}

func (tm *TextManager) GetTexture(name string) *sdl.Texture {
	if node, ok := tm.cache[name]; ok {
		return node.texture
	}
	return nil
}

func (tm *TextManager) createText(text string) (*sdl.Texture, error) {
	var solid *sdl.Surface
	var err error
	var texture *sdl.Texture
	white := sdl.Color{255, 255, 255, 255}
	// red := sdl.Color{255, 0, 0, 255}
	if solid, err = tm.font.RenderUTF8Blended(text, white); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write string: %s", text)
		return nil, err
	} else {
		texture, err = tm.engine.renderer.CreateTextureFromSurface(solid)
		defer solid.Free()
		if err != nil {
			fmt.Println("error converting text surface to texture")
		}
	}
	return texture, nil
}

func (tm *TextManager) Cleanup() {
	for _, v := range tm.cache {
		v.texture.Destroy()
	}
	tm.font.Close()
	ttf.Quit()
}
