package tips

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
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

func writeCSVToFile(path string, persons []person) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to create %s: %v", path, err)
	}
	defer func() {
		if cerr := f.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	w := csv.NewWriter(f)
	for _, p := range persons {
		row := []string{p.name, strconv.Itoa(p.age)}
		w.Write(row)
	}
	// Don't forget to call Flush!
	w.Flush()
	// Return w.Error(). Don't ignore an error occured in Flush().
	return w.Error()
}

func TestWriteCSV(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("testdata/tmp.csv"); err != nil {
			t.Error(err)
		}
	}()
	persons := []person{{name: "Taro", age: 8}, {name: "Jiro", age: 5}}
	err := writeCSVToFile("testdata/tmp.csv", persons)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile("testdata/tmp.csv")
	if err != nil {
		t.Error(err)
	}
	content := string(b)
	if content != "Taro,8\nJiro,5\n" {
		t.Errorf("Unexpected content: %q", content)
	}
}
