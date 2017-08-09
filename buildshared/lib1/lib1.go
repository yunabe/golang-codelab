package lib1

import (
	"log"
	"math"

	"github.com/yunabe/golang-codelab/buildshared/lib0"
)

var a, b float64 = 1.1, 2.0

func init() {
	log.Println("Package lib1 is initialized.")
	a = math.Pow(a, float64(lib0.GetX()))
	b = math.Pow(b, float64(lib0.GetY()))
	log.Printf("a = %f, b = %f", a, b)
}

func GetA() float64 {
	return a
}

func GetB() float64 {
	return b
}
