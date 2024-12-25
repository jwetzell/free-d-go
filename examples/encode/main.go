package main

import (
	"fmt"

	freeD "github.com/jwetzell/free-d-go"
)

func main() {

	freeDMessage := freeD.FreeDPosition{
		ID:    0xff,
		Pan:   -10.1,
		Tilt:  20.2,
		Roll:  3.33,
		PosX:  1.11,
		PosY:  2.22,
		PosZ:  3.33,
		Zoom:  1000,
		Focus: 2000,
	}

	freeDBytes := freeD.Encode(freeDMessage)

	println("send packet somehow")
	fmt.Printf("%v\n", freeDBytes)

}
