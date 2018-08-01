package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"fmt"
)

type Area struct {
	mapList []Map
	areaSize int32
	tileset *sdl.Texture
}

func (a *Area) Load(file string) error {
	fp, err := fileManager.GetArea(file)
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer fp.Close()

	var tileset string
	fmt.Fscanf(fp, "%s", &tileset)
	tileset = fileManager.GetTilesetPath(tileset)
	a.tileset, err = gfx.Load(game.renderer, tileset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load map tileset: %s\n", err)
		panic(err)
	}

	// get area size
	fmt.Fscanf(fp, "%d\n", &a.areaSize)

	for x := 0; int32(x) < a.areaSize; x++ {
		for y := 0; int32(y) < a.areaSize; y++ {
			var mapName string
			fmt.Fscanf(fp, "%s", &mapName)
			tmpMap := Map{}
			tmpMap.Load(mapName)
			tmpMap.Texture = a.tileset
			a.mapList = append(a.mapList, tmpMap)
		}
		fmt.Fscanf(fp, "\n")
	}
	return nil
}

func (a *Area) Render(cameraX int32, cameraY int32) {
	mapWidth := int32(MAP_WIDTH * TILE_SIZE)
	mapHeight := int32(MAP_HEIGHT * TILE_SIZE)

	firstId := -cameraX / mapWidth
	firstId = firstId + ((-cameraY / mapHeight) * a.areaSize)

	for i := 0; i < 4; i++ {
		id := firstId + ((int32(i) / 2) * a.areaSize) + (int32(i) % 2)

		if id < 0 || id >= int32(len(a.mapList)) {
			continue
		}

		x := ((id % a.areaSize) * mapWidth) + cameraX
		y := ((id / a.areaSize) * mapHeight) + cameraY

		a.mapList[id].Render(x, y)
	}
}

func (a *Area) GetMap(x int32, y int32) *Map {
	mapWidth := int32(MAP_WIDTH * TILE_SIZE)
	mapHeight := int32(MAP_HEIGHT * TILE_SIZE)
	id := x / mapWidth
	id = id + ((y / mapHeight) * a.areaSize)

	if id < 0 || id >= int32(len(a.mapList)) {
		return nil
	}

	return &a.mapList[id]
}

func (a *Area) GetTile(x int32, y int32) *Tile {
	mapWidth := int32(MAP_WIDTH * TILE_SIZE)
	mapHeight := int32(MAP_HEIGHT * TILE_SIZE)

	theMap := a.GetMap(x, y)

	if theMap == nil { return nil }

	x = x % mapWidth
	y = y % mapHeight

	return theMap.GetTile(x, y)
}

func (a *Area) Cleanup() {
	a.tileset.Destroy()
}
