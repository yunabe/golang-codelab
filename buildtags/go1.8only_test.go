// +build go1.8,!1.9

package buildtags

import (
	"runtime"
	"strings"
	"testing"
)

func TestGo1_8Only(t *testing.T) {
	version := runtime.Version()
	if !strings.HasPrefix(version, "go1.8") {
		t.Errorf("Unexpected version: %s", version)
	}
}
