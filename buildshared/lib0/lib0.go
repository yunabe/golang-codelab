package lib0

import (
	"log"
)

const x = 2

var y = 1

func init() {
	y = -1 * y
	log.Println("Package lib0 is initialized.")
}

func GetX() int {
	return x
}

func GetY() int {
	return y
}
