package types

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	_ "github.com/yunabe/easycsv" // Installs it in `go get`. This hack does not work from go1.10?
)

func TestTypesChecker(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "tmp.go", `
package tmp

import "github.com/yunabe/easycsv"

var r = easycsv.NewReader(nil)
`, /*mode*/ 0)
	if err != nil {
		t.Fatal(err)
	}
	pkg := types.NewPackage("github.com/yunabe/tmp", "tmp")
	config := &types.Config{
		Importer: importer.Default(),
	}
	info := &types.Info{
		Defs:   make(map[*ast.Ident]types.Object),
		Uses:   make(map[*ast.Ident]types.Object),
		Scopes: make(map[ast.Node]*types.Scope),
		Types:  make(map[ast.Expr]types.TypeAndValue),
	}
	ch := types.NewChecker(config, fset, pkg, info)
	if err := ch.Files([]*ast.File{f}); err != nil {
		t.Fatal(err)
	}
}
