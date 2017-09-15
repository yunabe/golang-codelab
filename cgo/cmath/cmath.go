package cmath

import (
	// #cgo LDFLAGS: -lm
	// #include <math.h>
	"C"
)

// Sin Wraps sin in <math.h>
func Sin(x float64) float64 {
	return float64(C.sin(C.double(x)))
}

// Cos Wraps sin in <math.h>
func Cos(x float64) float64 {
	return float64(C.cos(C.double(x)))
}
