package tips

import (
	"bufio"
	"io"
	"io/ioutil"
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
	// Use os.Open to open a file for read.
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

func TestReadUtils(t *testing.T) {
	func() {
		// ioutil.ReadFile read the entire file content as []byte.
		// There is no such function like ioutil.ReadFileString.
		bin, err := ioutil.ReadFile("testdata/example0.txt")
		if err != nil {
			t.Errorf("Failed to open a file with ioutil.ReadFile: %v", err)
		}
		str := string(bin)
		if str != "apple\nbanana \n cat\ndog\n" {
			t.Errorf("Unexpected data was read from example0.txt")
		}
	}()
	func() {
		// ioutil.ReadAll reads the entire content from io.Reader as []byte.
		f, err := os.Open("testdata/example0.txt")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Error(err)
			}
		}()
		bin, err := ioutil.ReadAll(f)
		if err != nil {
			t.Errorf("Failed to read a file with ioutil.ReadAll: %v", err)
			return
		}
		str := string(bin)
		if str != "apple\nbanana \n cat\ndog\n" {
			t.Errorf("Unexpected data was read from example0.txt")
		}
	}()
}

func removeTmpTxt(t *testing.T) {
	if err := os.RemoveAll("testdata/tmp.txt"); err != nil {
		t.Error(err)
	}
}

// An example to show how to use bufio.Writer.
// bufio.Write is Java's BufferedWriter in Go.
func TestBufioWriter(t *testing.T) {
	defer removeTmpTxt(t)

	// Use os.Create to open a file for write.
	f, err := os.Create("testdata/tmp.txt")
	if err != nil {
		t.Error(err)
		return
	}
	w := bufio.NewWriter(f)
	// It's okay to ignore the first return value (# of chars written) in WriteString.
	_, err = w.WriteString("Hello bufio.Writer.WriterString!")
	if err != nil {
		t.Error(err)
	}
	// It's okay to ignore the first return value in bufio.Writer.Write too.
	_, err = w.Write([]byte("Hello bufio.Writer.Writer!"))
	if err != nil {
		t.Error(err)
	}
	// DO NOT FORGET TO FLUSH bufio.Writer!
	// If you forget to call Flush, no data is actually written to the file.
	// Even worse, no runtime error is reported for that :{
	if err := w.Flush(); err != nil {
		t.Error(err)
	}

	// Check the content of the file.
	b, err := ioutil.ReadFile("testdata/tmp.txt")
	if err != nil {
		t.Error(err)
	}
	str := string(b)
	if str != "Hello bufio.Writer.WriterString!Hello bufio.Writer.Writer!" {
		t.Errorf("Unexpected content: %q", str)
	}
}

func TestWriteUtils(t *testing.T) {
	defer removeTmpTxt(t)

	err := ioutil.WriteFile("testdata/tmp.txt", []byte("Hello ioutil.WriteFile"), 0666)
	if err != nil {
		t.Error(err)
	}
	// Check the content of the file.
	b, err := ioutil.ReadFile("testdata/tmp.txt")
	if err != nil {
		t.Error(err)
	}
	str := string(b)
	if str != "Hello ioutil.WriteFile" {
		t.Errorf("Unexpected content: %q", str)
	}
}
