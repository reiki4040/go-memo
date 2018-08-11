package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	SeekFromFront()
	SeekFromFrontPlus()
	SeekFromEnd()
}

func SeekFromFront() {
	fp, err := os.Open("./data")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fp.Seek(5, 0)
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("offset 5, whence 0: %s\n", string(b))
}

func SeekFromFrontPlus() {
	fp, err := os.Open("./data")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fp.Seek(5, 0)
	fp.Seek(5, 1)
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("offset 5, whence 1: %s (after offset 5, whence 0)\n", string(b))
}

func SeekFromEnd() {
	fp, err := os.Open("./data")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fp.Seek(-3, 2)
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("offset -3, whence 2: %s\n", string(b))
}
