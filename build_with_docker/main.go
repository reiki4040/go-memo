package main

import (
	"fmt"

	"github.com/reiki4040/cstore"
)

func main() {
	// dummy outside third party library
	cstore.NewManager("", "")

	fmt.Println("Hello golang!")
}
