package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform/sdl"
	"github.com/bennicholls/tyumi/vec"
)

//import "github.com/pkg/profile"

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	log.EnableConsoleOutput()
	log.SetMinimumLogLevel(log.LVL_DEBUG)
	defer log.WriteToDisk()

	tyumi.SetPlatform(sdl.NewPlatform())
	tyumi.InitConsole("Spaceshippers", vec.Dims{96, 54}, "res/cp437_20x20.bmp", "res/DelveFont10x20.bmp")

	//console.SetFullscreen(true)
	//burl.Debug()

	startMenu := StartMenu{}
	startMenu.Init()
	tyumi.SetInitialMainState(&startMenu)

	tyumi.Run()
}
