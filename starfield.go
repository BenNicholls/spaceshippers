package main

import "math/rand"

func (sg *SpaceshipGame) initStarField() {
	w, h := sg.shipdisplay.Dims()
	sg.starField = make([]int, w*h)
	for i := 0; i < len(sg.starField); i++ {
		if rand.Intn(sg.starFrequency) == 0 {
			sg.starField[i] = 1
		}
	}
}

func (sg *SpaceshipGame) shiftStarField() {
	w, _ := sg.shipdisplay.Dims()
	for i := 0; i < len(sg.starField); i++ {
		if i%w != w-1 {
			sg.starField[i] = sg.starField[i+1]
		} else if rand.Intn(sg.starFrequency) == 0 {
			sg.starField[i] = 1
		} else {
			sg.starField[i] = 0
		}
	}
}

func (sg *SpaceshipGame) DrawStarfield() {
	sg.shipdisplay.Clear()
	w, h := sg.shipdisplay.Dims()
	for i := 0; i < w*h; i++ {
		if sg.starField[i] != 0 {
			sg.shipdisplay.Draw(i%w, i/w, 0x2a, 0xFF444444, 0xFF000000)
		}
	}
}
