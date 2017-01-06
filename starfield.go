package main

import "math/rand"

var starField []int

const starFrequency int = 20

func initStarField() {
	for i := 0; i < len(starField); i++ {
		if rand.Intn(starFrequency) == 0 {
			starField[i] = 1
		}
	}
}

func shiftStarField() {
	w, _ := shipdisplay.Dims()
	for i := 0; i < len(starField); i++ {
		if i%w != w-1 {
			starField[i] = starField[i+1]
		} else if rand.Intn(starFrequency) == 0 {
			starField[i] = 1
		} else {
			starField[i] = 0
		}
	}
}

func DrawStarfield() {
	shipdisplay.Clear()
	w, h := shipdisplay.Dims()
	for i := 0; i < w*h; i++ {
		if starField[i] != 0 {
			shipdisplay.Draw(i%w, i/w, 0x2a, 0xFF444444, 0xFF000000)
		}
	}
}
