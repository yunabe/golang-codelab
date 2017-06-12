package tips

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func readCSVFromFile(path string) (persons []person, err error) {
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
		age, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %q as an integer", row[1])
		}
		persons = append(persons, person{name: row[0], age: int(age)})
	}
	return persons, nil
}

func TestReadCSV(t *testing.T) {
	persons, err := readCSVFromFile("testdata/example.csv")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(persons, []person{
		{name: "Alice", age: 20},
		{name: "Bob", age: 23},
	}) {
		t.Error("Unexpected result: ", persons)
	}
}
