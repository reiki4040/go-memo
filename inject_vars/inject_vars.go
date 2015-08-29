package main

import (
	"fmt"
)

var (
	version   string
	hash      string
	builddate string
	goversion string
)

func main() {
	fmt.Printf("version: %s (%s)\n", version, hash)
	fmt.Printf("build at %s with %s\n", builddate, goversion)
}
