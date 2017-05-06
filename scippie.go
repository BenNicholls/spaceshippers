package main

import "strings"

func (sg *SpaceshipGame) Execute() {
	sg.output.Append("")
	sg.output.Append(">>> " + sg.input.GetText())
	sg.output.Append("")
	switch strings.ToLower(sg.input.GetText()) {
	case "status":
		for _, r := range sg.PlayerShip.Rooms {
			sg.output.Append(r.GetStatus())
		}
	case "help":
		sg.output.Append("S.C.I.P.P.I.E. is your AI helper. Give him one of the following commands, and he'll get 'r done!")
		sg.output.Append("   status     prints ship room status")
		sg.output.Append("   help       prints a mysterious menu")
	default:
		sg.output.Append("I do not understand that command, you dummo. Try \"help\"")
	}
	sg.output.ScrollToBottom()
}

func (sg *SpaceshipGame) AddMessage(s string) {
	sg.output.Append(s)
	sg.output.ScrollToBottom()
}