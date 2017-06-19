package easycsv

import (
	"fmt"
	"log"
	"os"
)

func ExampleReader_read() {
	f, err := os.Open("testdata/sample.csv")
	if err != nil {
		log.Fatalf("Failed to open a file: %v", err)
	}
	r := NewReadCloser(f)
	var entry struct {
		Name string `index:"0"`
		Age  int    `index:"1"`
	}
	for r.Read(&entry) {
		fmt.Print(entry)
	}
	if err := r.Done(); err != nil {
		log.Fatalf("Failed to read a CSV file: %v", err)
	}
	// Output: {Alice 10}{Bob 20}
}
