// A very simple tool to demonstrate the usage of go/build API.
package main

import (
	"flag"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var pkgPath = flag.String("package_path", "github.com/yunabe/golang-codelab/misc/analyzedeps", "The path of a package to analyze")

func isStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}

func main() {
	flag.Parse()

	pkg, err := build.Default.Import(*pkgPath, filepath.Join(runtime.GOROOT(), "src"), 0)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("import path == %s (is standard == %t)\n", pkg.ImportPath, isStandardImportPath(pkg.ImportPath))
	log.Printf("imports == %#v", pkg.Imports)
}
