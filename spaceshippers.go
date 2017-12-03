package main

import "fmt"
import "math/rand"
import "time"
import "github.com/bennicholls/burl-E/burl"

//import "github.com/pkg/profile"

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	rand.Seed(time.Now().UTC().UnixNano())

	console, err := burl.InitConsole(80, 45, "res/curses24x24.bmp", "res/DelveFont12x24.bmp", "Spaceshippers")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	console.SetFullscreen()

	ssg := NewSpaceshipGame()

	burl.InitState(ssg)
	err = burl.GameLoop()

	if err != nil {
		fmt.Println(err)
		return
	}

	ssg.SaveShip()
}
