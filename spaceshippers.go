package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bennicholls/burl"
	"github.com/bennicholls/burl/console"
)

func main() {

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
