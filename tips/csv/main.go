package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	inputPath = flag.String("input", "", "Input Path")
)

// Person holds the content of CSV.
type Person struct {
	Name string
	Age  int32
}

func readCSV(path string) (persons []Person, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %v", path, err)
	}
	defer func() {
		if cerr := f.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	r := csv.NewReader(f)
	for lineno := 1; ; lineno++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if l := len(row); l != 2 {
			return nil, fmt.Errorf("The row size must be 2 but got %d at line %d", l, lineno)
		}
		age, err := strconv.ParseInt(row[1], 0, 32)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %q as an integer", row[1])
		}
		persons = append(persons, Person{Name: row[0], Age: int32(age)})
	}
	return persons, nil
}

func main() {
	flag.Parse()

	_, err := readCSV(*inputPath)
	if err != nil {
		log.Fatal(err)
	}
}
