package main

import (
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type MainMenu struct {
	ui.PageContainer
}

func (mm *MainMenu) Init() {
	mm.PageContainer.Init(vec.Dims{56, 45}, vec.Coord{39, 4}, 10)
	mm.PageContainer.EnableBorder()
	mm.Hide()

	return
}