package main

import (
	"log"
	"unsafe"

	"github.com/yunabe/golang-codelab/buildshared/lib0"
)

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
*/
import "C"

func main() {
	log.Printf("lib0.GetX() == %d, lib0.GetY() == %d\n", lib0.GetX(), lib0.GetY())

	// I wrote this code by referring https://golang.org/src/plugin/plugin_dlopen.go
	handle := C.dlopen(C.CString("cache/libgithub.com-yunabe-golang-codelab-buildshared-lib1.so"), C.RTLD_LAZY)
	if handle == nil {
		panic("Failed to open shared object.")
	}
	initFuncPC := C.dlsym(handle, C.CString("github.com/yunabe/golang-codelab/buildshared/lib1.init"))
	if initFuncPC == nil {
		panic("Could not found function pointer.")
	}
	initFuncP := &initFuncPC
	initFunc := *(*func())(unsafe.Pointer(&initFuncP))
	initFunc()
}
