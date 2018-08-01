package main

const (
	TILE_TYPE_NONE = 0
	TILE_TYPE_NORMAL = 1
	TILE_TYPE_BLOCK = 2
)

type Tile struct {
	TileID int32
	TypeID int32
}