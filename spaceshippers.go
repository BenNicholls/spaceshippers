package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/bennicholls/burl/console"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	runtime.LockOSThread()
	rand.Seed(time.Now().UTC().UnixNano())

	err := console.Setup(80, 45, "res/curses24x24.bmp", "res/DelveFont12x24.bmp", "Spaceshippers")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	console.SetFullscreen()

	var event sdl.Event
	running := true

	m := NewSpaceshipGame()

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.WindowEvent:
			 	if t.Event == sdl.WINDOWEVENT_RESTORED {
					 console.ForceRedraw()
				 }
			// case *sdl.MouseMotionEvent:
			// 	fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			// case *sdl.MouseButtonEvent:
			// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			// case *sdl.MouseWheelEvent:
			// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				m.HandleKeypress(t.Keysym.Sym)
			}
		}

		m.Update()

		m.Render()
		console.Render()
	}
}


