package tips

import (
	"testing"

	// The leaf directory name does not necessarily match to the package default name.
	"github.com/yunabe/golang-codelab/tips/importpath/my-pkg"

	"github.com/yunabe/golang-codelab/tips/importpath/my_pkg"
	"github.com/yunabe/golang-codelab/tips/importpath/mypkg"
)

func TestImportPathName(t *testing.T) {
	if mypkg.Var != 30 {
		t.Errorf("Unexpected mypkg.Var: %d", mypkg.Var)
	}
	if my_pkg.Var != 10 {
		t.Errorf("Unexpected my_pkg.Var: %d", my_pkg.Var)
	}
	if mydashpkg.Var != 20 {
		t.Errorf("Unexpected mydashpkg.Var: %d", mydashpkg.Var)
	}
}
