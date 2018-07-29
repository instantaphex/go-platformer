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
	fp, err := os.Open(file)

	if err != nil {
		fmt.Print(err)
		return err
	}
	defer fp.Close()

	// grab tileset from area and load it
	var tilesetFile string
	cwd, _ := os.Getwd()
	fmt.Fscanf(fp, "%s", &tilesetFile)
	tilesetFile = cwd + tilesetFile
	a.tileset, err = gfx.Load(game.renderer, tilesetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load map tileset: %s\n", err)
		panic(err)
	}

	// get area size
	fmt.Fscanf(fp, "%d\n", &a.areaSize)

	for x := 0; int32(x) < a.areaSize; x++ {
		for y := 0; int32(y) < a.areaSize; y++ {
			var mapFile string
			fmt.Fscanf(fp, "%s", &mapFile)
			cwd, _ := os.Getwd()
			tmpMap := Map{}
			tmpMap.Load(cwd + mapFile)
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

func (a *Area) Cleanup() {
	a.tileset.Destroy()
}
