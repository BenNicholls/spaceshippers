package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform/sdl"
	"github.com/bennicholls/tyumi/vec"
	//"github.com/pkg/profile"
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	defer log.WriteToDisk()

	tyumi.SetPlatform(sdl.NewPlatform())
	tyumi.InitConsole("Spaceshippers", vec.Dims{96, 54}, "res/cp437_20x20.bmp", "res/DelveFont10x20.bmp")

	ui.SetDefaultFocusColour(col.ORANGE)
	tyumi.EnableCursor()

	//console.SetFullscreen(true)

	startMenu := StartMenu{}
	startMenu.Init()
	tyumi.SetInitialScene(&startMenu)

	tyumi.Run()
}
