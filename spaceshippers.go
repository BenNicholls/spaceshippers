package main

import "fmt"
import "math/rand"
import "time"
import "github.com/bennicholls/burl"
import "github.com/bennicholls/burl/console"
//import "github.com/pkg/profile"

func main() {
	//defer profile.Start(profile.ProfilePath(".")).Stop()
	rand.Seed(time.Now().UTC().UnixNano())

	err := console.Setup(80, 45, "res/curses24x24.bmp", "res/DelveFont12x24.bmp", "Spaceshippers")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	console.SetFullscreen()

	burl.InitState(NewSpaceshipGame())
	err = burl.GameLoop()

	if err != nil {
		fmt.Println(err)
		return
	}
}
