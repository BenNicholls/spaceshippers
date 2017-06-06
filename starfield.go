package main

import "math/rand"
import "github.com/bennicholls/burl/util"

//initializes a starfield twice the width of the screen
func (sg *SpaceshipGame) initStarField(starFrequency int) {
	w, h := sg.shipdisplay.Dims()
	sg.starField = make([]int, w*h*2)
	for i := 0; i < len(sg.starField); i++ {
		if rand.Intn(starFrequency) == 0 {
			sg.starField[i] = 1
		}
	}
}

//moves the "camera" on the stars.
func (sg *SpaceshipGame) shiftStarField() {
	w, _ := sg.shipdisplay.Dims()
	sg.starShift, _ = util.ModularClamp(sg.starShift+1, 0, (w*2)-1)
}

//Draws the starfield, offset by the current starShift value.
func (sg *SpaceshipGame) DrawStarfield() {
	sg.shipdisplay.Clear()
	w, h := sg.shipdisplay.Dims()
	for i := 0; i < w*h; i++ {
		if sg.starField[(i/w)*w*2+(i%w+sg.starShift)%(w*2)] != 0 {
			sg.shipdisplay.Draw(i%w, i/w, 0x2a, 0xFF444444, 0xFF000000)
		}
	}
}
