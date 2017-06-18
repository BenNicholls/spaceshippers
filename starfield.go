package main

import "math/rand"
import "github.com/bennicholls/burl/util"
import "github.com/bennicholls/burl/ui"

type StarField struct {
	field         []int
	starFrequency int
	starShift     int
	view          *ui.TileView
	dirty         bool
}

//initializes a starfield twice the width of the screen
func NewStarField(w, h, starFrequency int, v *ui.TileView) (sf StarField) {
	sf.view = v
	sf.field = make([]int, w*h*2)
	sf.starFrequency = starFrequency
	sf.dirty = true
	for i := 0; i < len(sf.field); i++ {
		if rand.Intn(sf.starFrequency) == 0 {
			sf.field[i] = 1
		}
	}
	return
}

//moves the "camera" on the stars.
func (sf *StarField) Shift() {
	w, _ := sf.view.Dims()
	sf.starShift, _ = util.ModularClamp(sf.starShift+1, 0, (w*2)-1)
	sf.dirty = true
}

//Draws the starfield, offset by the current starShift value.
func (sf *StarField) Draw() {
	if sf.dirty {
		sf.view.Reset()
		w, h := sf.view.Dims()
		for i := 0; i < w*h; i++ {
			if sf.field[(i/w)*w*2+(i%w+sf.starShift)%(w*2)] != 0 {
				sf.view.Draw(i%w, i/w, 0x2a, 0xFF444444, 0xFF000000)
			}
		}
		sf.dirty = false
	}
}
