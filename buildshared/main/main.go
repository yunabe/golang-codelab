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

	libs := []struct {
		path string
		init string
	}{
		{
			path: "cache/libgithub.com-yunabe-golang-codelab-buildshared-lib1.so",
			init: "github.com/yunabe/golang-codelab/buildshared/lib1.init",
		},
		{
			path: "cache/libgithub.com-yunabe-golang-codelab-buildshared-lib3.so",
			init: "github.com/yunabe/golang-codelab/buildshared/lib3.init",
		},
	}
	for _, lib := range libs {
		// I wrote this code by referring https://golang.org/src/plugin/plugin_dlopen.go
		handle := C.dlopen(C.CString(lib.path), C.RTLD_LAZY)
		if handle == nil {
			panic("Failed to open shared object.")
		}
		initFuncPC := C.dlsym(handle, C.CString(lib.init))
		if initFuncPC == nil {
			panic("Could not found function pointer.")
		}
		initFuncP := &initFuncPC
		initFunc := *(*func())(unsafe.Pointer(&initFuncP))
		initFunc()
	}
}
