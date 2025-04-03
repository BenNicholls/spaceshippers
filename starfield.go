package main

import (
	"math/rand"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type StarField struct {
	ui.Element

	field          []bool
	offset         int
	shiftFrequency int // number of game ticks between shifts. if 0, no shift.
}

func (sf *StarField) Init(size vec.Dims, pos vec.Coord, depth, starFrequency, shiftFrequency int) {
	sf.Element.Init(size, pos, depth)
	sf.TreeNode.Init(sf)
	sf.shiftFrequency = shiftFrequency

	sf.field = make([]bool, size.Area()*2)
	for i := range sf.field {
		if rand.Intn(starFrequency) == 0 {
			sf.field[i] = true
		}
	}
}

func (sf *StarField) Update() {
	if sf.shiftFrequency == 0 {
		return
	}
	
	// moves the "camera" on the stars.
	if tyumi.GetTick()%sf.shiftFrequency != 0 {
		return
	}

	sf.offset = util.CycleClamp(sf.offset+1, 0, (sf.Size().W*2)-1)
	sf.Updated = true
}

// Draws the starfield, offset by the current starShift value.
func (sf *StarField) Render() {
	w := sf.Size().W
	star := gfx.NewGlyphVisuals(gfx.GLYPH_ASTERISK, col.Pair{col.DARKGREY, col.NONE})
	for cursor := range vec.EachCoordInArea(sf.DrawableArea()) {
		if sf.field[(cursor.Y)*w*2+(cursor.X+sf.offset)%(w*2)] {
			sf.DrawVisuals(cursor, 0, star)
		} else {
			sf.DrawGlyph(cursor, 0, gfx.GLYPH_SPACE)
		}
	}
}
