package main

import "fmt"
import "math/rand"
import "time"
import "github.com/bennicholls/burl-E/burl"

//import "github.com/pkg/profile"

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	rand.Seed(time.Now().UTC().UnixNano())

	_, err := burl.InitConsole(96, 54, "res/cp437_20x20.bmp", "res/DelveFont10x20.bmp", "Spaceshippers")
	if err != nil {
		fmt.Println(err)
		return
	}

	burl.Debug()

	burl.InitState(NewStartMenu())

	err = burl.GameLoop()
	if err != nil {
		fmt.Println(err)
	}
}
