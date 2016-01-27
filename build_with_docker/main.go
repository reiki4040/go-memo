package main

import (
	"fmt"

	"github.com/reiki4040/cstore"
)

var (
	version   string
	hash      string
	goversion string
)

func main() {
	// dummy outside third party library
	cstore.NewManager("", "")

	fmt.Println("Hello golang!!")
	fmt.Printf("version: %s (%s) golang %s\n", version, hash, goversion)
}
