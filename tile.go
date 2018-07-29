package main

const (
	TILE_TYPE_NONE = 0 << iota
	TILE_TYPE_NORMAL
	TILE_TYPE_BLOCK
)

type Tile struct {
	TileID int32
	TypeID int32
}