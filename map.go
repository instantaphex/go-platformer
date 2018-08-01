package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"io/ioutil"
)

type Map struct {
	Texture *sdl.Texture
	Height int32
	Width int32
	TileSize int32
	tileList []Tile
}

/*func (m *Map) Load(mapName string) bool {
	f, err := fileManager.GetMap(mapName)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer f.Close()

	for y := 0; y < m.Height; y++ {
		for x:= 0; x < m.Width; x++ {
			tmpTile := Tile{}
			fmt.Fscanf(f, "%d:%d", &tmpTile.TileID, &tmpTile.TypeID)
			m.tileList = append(m.tileList, tmpTile)
		}
		fmt.Fscanf(f, "\n")
	}
	return true
}*/

func (m *Map) Load(mapName string) bool {
	var err error

	tmx := m.openTmx(mapName)
	tilePath := fileManager.GetTilesetPath("tilemap")
	m.Texture, err = gfx.Load(game.renderer, tilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load tileset: %s\n", err)
		return false
	}

	m.Height = int32(tmx.Height)
	m.Width = int32(tmx.Width)
	m.TileSize = int32(tmx.TileWidth)

	for _, v := range tmx.Layers[0].Data.ParsedData {
		var tileId, typeId int32
		if v == 0 {
			tileId = 33
			typeId = TILE_TYPE_NONE
		} else {
			tileId = int32(v - 1)
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

	tilesetWidth := w / m.TileSize
	tilesetHeight := h / m.TileSize

	id := 0


	for y := int32(0); y < m.Height; y++ {
		for x:= int32(0); x < m.Width; x++ {
			if m.tileList[id].TypeID == TILE_TYPE_NONE {
				id++
				continue
			}

			tX := mapX + int32(x * m.TileSize)
			tY := mapY + int32(y * m.TileSize)

			tilesetX := (m.tileList[id].TileID % tilesetWidth) * m.TileSize
			tilesetY := (m.tileList[id].TileID / tilesetHeight) * m.TileSize

			gfx.DrawPart(game.renderer, m.Texture, tX, tY, tilesetX, tilesetY, m.TileSize, m.TileSize)

			id++
		}
	}
}

func (m *Map) GetTile(x int32, y int32) *Tile {
	id := x / m.TileSize
	id = id + (m.Width * (y / m.TileSize))
	if id < 0 || id >= int32(len(m.tileList)) {
		return nil
	}
	return &m.tileList[id]
}

func (m *Map) Cleanup() {
	m.Texture.Destroy()
}

func (m Map) openTmx(filename string) TmxMap {
	path := fileManager.GetPath("tmx", filename, "tmx")
	f, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	parsed, _ := Parse(f)
	return parsed
}
