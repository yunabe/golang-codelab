# buildshared
Files in this directory demostrate how to build shared libraries in Go and how to load then dynamically.
- How to build Go shared libraries efficiently without building all source files from scratch.
- How to dynamically load Go shared libraries.

# Commands
``` shell
go install -buildmode=shared -pkgdir `pwd`/cache -linkshared github.com/yunabe/golang-codelab/buildshared/lib0
go build -pkgdir `pwd`/cache -linkshared -o buildshareddemo github.com/yunabe/golang-codelab/buildshared/main
ldd ./buildshareddemo
```

## Notes
- The path for `-pkgdir` must be an absolute PATH.
