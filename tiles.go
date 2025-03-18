package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl"
)

var TILE_FLOOR = rl.RegisterTileType(rl.TileData{
	Name:     "Floor",
	Passable: true,
	Glyph:    gfx.GLYPH_DOT_SMALL,
	Colours:  col.Pair{col.DARKGREY, col.NONE},
})

var TILE_WALL = rl.RegisterTileType(rl.TileData{
	Name:    "Wall",
	Opaque:  true,
	Glyph:   gfx.GLYPH_HASH,
	Colours: col.Pair{col.GREY, col.NONE},
})

var TILE_DOOR = rl.RegisterTileType(rl.TileData{
	Name:     "Door",
	Passable: true,
	Opaque:   true,
	Glyph:    gfx.GLYPH_IDENTICAL,
	Colours:  col.Pair{col.GREY, col.NONE},
})
