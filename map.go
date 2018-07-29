package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

type Map struct {
	Texture *sdl.Texture
	tileList []Tile
}

func (m *Map) Load(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer f.Close()

	for y := 0; y < MAP_HEIGHT; y++ {
		for x:= 0; x < MAP_WIDTH; x++ {
			tmpTile := Tile{}
			fmt.Fscanf(f, "%d:%d", &tmpTile.TileID, &tmpTile.TypeID)
			m.tileList = append(m.tileList, tmpTile)
		}
		fmt.Fscanf(f, "\n")
	}
	return true
}

func (m *Map) Render(mapX int32, mapY int32) {
	if m.Texture == nil {
		return
	}

	_, _, w, h, err := m.Texture.Query()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	tilesetWidth := w / TILE_SIZE
	tilesetHeight := h / TILE_SIZE

	id := 0


	for y := 0; y < MAP_HEIGHT; y++ {
		for x:= 0; x < MAP_WIDTH; x++ {
			if m.tileList[id].TypeID == TILE_TYPE_NONE {
				id++
				continue
			}

			tX := mapX + int32(x * TILE_SIZE)
			tY := mapY + int32(y * TILE_SIZE)

			tilesetX := (m.tileList[id].TileID % tilesetWidth) * TILE_SIZE
			tilesetY := (m.tileList[id].TileID / tilesetHeight) * TILE_SIZE

			gfx.DrawPart(game.renderer, m.Texture, tX, tY, tilesetX, tilesetY, TILE_SIZE, TILE_SIZE)

			id++
		}
	}
}
