package tips

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"testing"
)

// readLinesByScanner shows how to read a file line by line in Go with bufio.Scanner.
// Notes:
// - Unlike readlines in Python, Scanner.Text() returns a text without '\n'.
// - Scanner ignores the last line of a file if the last line is empty.
// - You do not need to wrap f with bufio.Reader because bufio.Scanner also has
//   an internal buffer to read the file content efficiently.
func readLinesByScanner(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := f.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func TestReadLinesByScanner(t *testing.T) {
	lines, err := readLinesByScanner("testdata/example0.txt")
	if err != nil {
		t.Error(err)
	} else {
		expect := []string{"apple", "banana ", " cat", "dog"}
		if !reflect.DeepEqual(expect, lines) {
			t.Errorf("Expected %v but got %v", expect, lines)
		}
	}

	lines, err = readLinesByScanner("testdata/example1.txt")
	if err != nil {
		t.Error(err)
	} else {
		expect := []string{"alpha", "beta"}
		if !reflect.DeepEqual(expect, lines) {
			t.Errorf("Expected %v but got %v", expect, lines)
		}
	}
}

// You can use bufio.Reader to read lines though using bufio.Scanner is highly encouraged.
// - ReadString does not strip '\n'.
// - ReadString deos not ignore the last line even if it's empty.
// - Do not ignore the returned text when io.EOF is returned.
func readLinesByReader(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := f.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	r := bufio.NewReader(f)
	for {
		// line includes '\n'.
		line, err := r.ReadString('\n')
		if err == nil {
			lines = append(lines, line)
		} else if err == io.EOF {
			// BE CAREFUL! Don't ignore the text when io.EOF is returned.
			// Otherwise, you would lose the last line if it does not end with '\n'.
			lines = append(lines, line)
			break
		} else {
			return nil, err
		}
	}
	return lines, nil
}

func TestReadLinesByReader(t *testing.T) {
	lines, err := readLinesByReader("testdata/example0.txt")
	if err != nil {
		t.Error(err)
	} else {
		expect := []string{"apple\n", "banana \n", " cat\n", "dog\n", ""}
		if !reflect.DeepEqual(expect, lines) {
			t.Errorf("Expected %v but got %v", expect, lines)
		}
	}

	lines, err = readLinesByReader("testdata/example1.txt")
	if err != nil {
		t.Error(err)
	} else {
		expect := []string{"alpha\n", "beta"}
		if !reflect.DeepEqual(expect, lines) {
			t.Errorf("Expected %v but got %v", expect, lines)
		}
	}
}
