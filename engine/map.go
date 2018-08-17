package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
)

type Map struct {
	Texture *sdl.Texture
	Height int32
	Width int32
	TileSize int32
	tileList []Tile
	engine *Engine
}

const (
	TILE_TYPE_NONE = 0
	TILE_TYPE_NORMAL = 1
	TILE_TYPE_BLOCK = 2
)

type Tile struct {
	TileID int32
	TypeID int32
}

func (m *Map) Load(mapName string) error {
	var err error
	path := m.engine.File.GetPath("tmx", mapName, "tmx")
	tmx, err := NewTmxMap(path) // m.openTmx(mapName)
	if err != nil {
		return err
	}

	tilePath := m.engine.File.GetTilesetPath("tilemap")
	m.Texture, err = m.engine.Graphics.Load(tilePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load tileset: %s\n", err)
		return err
	}

	m.Height = int32(tmx.Height)
	m.Width = int32(tmx.Width)
	m.TileSize = int32(tmx.TileWidth)

	layer, err := tmx.GetLayerByName("map")
	if err != nil {
		fmt.Println(err)
	}

	group, err := tmx.GetObjGroupByName("objects")
	if err != nil {
		fmt.Println(err)
		fmt.Println(group)
	}
	/*for _, obj := range group.Objects {
		if obj.Type == "coin" {
			EntityList = append(EntityList, NewCoin(int32(obj.X), int32(obj.Y)))
		}
		if obj.Type == "heart" {
			EntityList = append(EntityList, NewHeart(int32(obj.X), int32(obj.Y)))
		}
	}*/

	for _, v := range layer.Data.ParsedData {
		var tileId, typeId int32
		if v == 0 {
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
	return nil
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

			m.engine.Graphics.DrawPart(m.Texture, tX, tY, tilesetX, tilesetY, m.TileSize, m.TileSize)

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


