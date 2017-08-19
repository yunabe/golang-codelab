package tips

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestCaptureStdout(t *testing.T) {
	// This test demostrates how to use pipe, dup2 syscalls to capture Stdout in Go.
	var rw [2]int
	err := syscall.Pipe(rw[:])
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := syscall.Close(rw[0]); err != nil {
			t.Error(err)
		}
		if err := syscall.Close(rw[1]); err != nil {
			t.Error(err)
		}
	}()
	nfd, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := syscall.Close(nfd); err != nil {
			t.Error(err)
		}
	}()
	syscall.Dup2(rw[1], syscall.Stdout)

	// c.f. os.Stdin
	done := make(chan struct{})
	msg := "Hello World!"
	var buf = make([]byte, len(msg))
	go func() {
		defer close(done)
		r := os.NewFile(uintptr(rw[0]), "pipereader")
		n, err := r.Read(buf)
		if err != nil {
			t.Error(err)
		}
		if n != len(msg) {
			t.Errorf("n is too small: %d < %d", n, len(msg))
		}
	}()
	fmt.Print(msg)
	syscall.Dup2(nfd, syscall.Stdout)
	<-done

	if string(buf) != msg {
		t.Errorf("Expected %q but got %q", msg, buf)
	}
}
