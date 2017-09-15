package main

import "fmt"
import "github.com/yunabe/golang-codelab/cgo/cmath"

func main() {
	fmt.Println(cmath.Sin(0))
	fmt.Println(cmath.Cos(0))
}
