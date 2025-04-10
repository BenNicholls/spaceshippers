package main

import (
	"github.com/bennicholls/tyumi/gfx/ui"
)

type MainMenu struct {
	ui.PageContainer
}

func (mm *MainMenu) Init() {
	mm.PageContainer.Init(menuSize, menuPos, menuDepth)
	mm.PageContainer.EnableBorder()
	mm.Hide()

	return
}
