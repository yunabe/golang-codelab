package main

import (
	// #include <math.h>
	"C"
	"fmt"
)

func main() {
	fmt.Println(C.cos(0))
	fmt.Println(C.sin(0))
}
