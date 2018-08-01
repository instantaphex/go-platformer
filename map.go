package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
)

type Map struct {
	Texture *sdl.Texture
	tileList []Tile
}

/*func (m *Map) Load(mapName string) bool {
	f, err := fileManager.GetMap(mapName)
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
}*/

func (m *Map) Load(mapName string) bool {
	tmx := openTMX(mapName)
	for _, v := range tmx.Layers[0].Data.ParsedData {
		var tileId, typeId int32
		if v == 0 {
			tileId = 33
			typeId = TILE_TYPE_NONE
		} else {
			tileId = int32(v)
			typeId = TILE_TYPE_BLOCK
		}
		tmpTile := Tile{}
		tmpTile.TileID = tileId
		tmpTile.TypeID = typeId
		m.tileList = append(m.tileList, tmpTile)
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

func (m *Map) GetTile(x int32, y int32) *Tile {
	id := x / TILE_SIZE
	id = id + (MAP_WIDTH * (y / TILE_SIZE))
	if id < 0 || id >= int32(len(m.tileList)) {
		return nil
	}
	return &m.tileList[id]
}
