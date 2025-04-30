package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/rl"
)

var TILE_FLOOR = rl.RegisterTileType(rl.TileData{
	Name:     "Floor",
	Passable: true,
	Visuals:  gfx.NewGlyphVisuals(gfx.GLYPH_DOT_SMALL, col.Pair{col.DARKGREY, col.NONE}),
})

var TILE_WALL = rl.RegisterTileType(rl.TileData{
	Name:    "Wall",
	Opaque:  true,
	Visuals: gfx.NewGlyphVisuals(gfx.GLYPH_HASH, col.Pair{col.GREY, col.NONE}),
})

var TILE_DOOR = rl.RegisterTileType(rl.TileData{
	Name:     "Door",
	Passable: true,
	Opaque:   true,
	Visuals:  gfx.NewGlyphVisuals(gfx.GLYPH_IDENTICAL, col.Pair{col.GREY, col.NONE}),
})
