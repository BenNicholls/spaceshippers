package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/delvetown/console"
import "github.com/bennicholls/delvetown/ui"
import "github.com/bennicholls/delvetown/util"
import "runtime"
import "fmt"

var window *ui.Container
var input *ui.Inputbox
var output *ui.List
var crew *ui.List
var status *ui.Container

func main() {

    runtime.LockOSThread()

    var event sdl.Event

    err := console.Setup(96, 54, "curses.bmp", "Spaceshippers")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer console.Cleanup()

    SetupUI()

    running := true

    for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
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
				HandleKeypress(t.Keysym.Sym)
			}
		}

        Update()

        Render()
        console.Render()
    }
}

func SetupUI() {
    window = ui.NewContainer(94, 52, 1, 1, 0, true)
    window.SetTitle("SPACESHIPPERS. THE GAME OF SPACE SHIPS.")
    window.ToggleFocus()
    input = ui.NewInputbox(70, 1, 1, 50, 0, true)
    input.ToggleFocus()
    input.SetTitle("SCIPP V6.18")
    output = ui.NewList(70, 47, 1, 1, 0, true, "Ask SCIPP a question, or give him a command!")
    output.ToggleHighlight()
    // var crew *ui.List
    // var status *ui.Container

    window.Add(input, output)
}

func HandleKeypress(key sdl.Keycode) {
    if util.ValidText(rune(key)) {
        input.InsertText(rune(key))
    } else {
        switch key {
        case sdl.K_RETURN:
            Execute()
            input.Reset()
        case sdl.K_BACKSPACE:
            input.Delete()
        case sdl.K_SPACE:
            input.Insert(" ")
        case sdl.K_UP:
            output.ScrollUp()
        case sdl.K_DOWN:
            output.ScrollDown()
        }
    }
}

func Update() {

}

func Render() {
    window.Render()
}

func Execute() {
    output.Append("")
    output.Append(">>> " + input.GetText())
    output.Append("")
    output.Append("I do not understand that command, you dummo.")
    output.Select(len(output.Elements) - 1)
}