package main

import (
	"fmt"

	"github.com/mgutz/ansi"
)

func main() {
	reset := ansi.ColorCode("reset")
	msg := ansi.Color("red+b", "red+b:white")
	fmt.Println(msg, reset)
	msg = ansi.Color("green", "green")
	fmt.Println(msg, reset)
	msg = ansi.Color("background white", "black:white")
	fmt.Println(msg, reset)

	warn := ansi.ColorFunc("yellow:black")
	fmt.Println(warn("this is warning!"), reset)

	lime := ansi.ColorCode("green+h:black")

	fmt.Println(lime, "Lime message.", reset)
}
