package goroutineid

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

var goroutineSpace = []byte("goroutine ")

// GetGoroutineID extracts the ID of the current goroutine from runtime.Stack().
// This function is implemented based on https://github.com/jtolds/gls/blob/master/gid.go
func GetGoroutineID() uint64 {
	var buf [64]byte
	b := buf[:runtime.Stack(buf[:], false)]

	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}
