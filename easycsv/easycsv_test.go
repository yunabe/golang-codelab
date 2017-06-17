package easycsv

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)
import "bytes"

func TestLoopNil(t *testing.T) {
	f := bytes.NewReader([]byte(""))
	r := NewReader(f)
	r.Loop(nil)
	err := r.Done()
	if err == nil || !strings.Contains(err.Error(), "must not be nil") {
		t.Errorf("Unexpected error: %v", err)
	}
}

type fakeCloser struct {
	err    error
	closed bool
}

func (c *fakeCloser) Close() error {
	c.closed = true
	return c.err
}

func (*fakeCloser) Read([]byte) (int, error) {
	return 0, nil
}

func TestCloser(t *testing.T) {
	c := &fakeCloser{}
	r := NewReadCloser(c)
	if err := r.Done(); err != nil {
		t.Error(err)
		return
	}
	if !c.closed {
		t.Error("c is not closed.")
	}
}

func TestCloserWithError(t *testing.T) {
	c := &fakeCloser{}
	c.err = errors.New("Close Error")
	r := NewReadCloser(c)
	if err := r.Done(); err != c.err {
		t.Errorf("Unexpected error: %v", err)
	}
	if !c.closed {
		t.Error("c is not closed.")
	}
}

func TestCloserDontOverwriteError(t *testing.T) {
	c := &fakeCloser{}
	c.err = errors.New("Close Error")
	r := NewReadCloser(c)
	anotherErr := errors.New("Another error")
	r.err = anotherErr
	if err := r.Done(); err != anotherErr {
		t.Errorf("Unexpected error: %v", err)
	}
	if !c.closed {
		t.Error("c is not closed.")
	}
}

func TestLoop(t *testing.T) {
	f := bytes.NewReader([]byte("10,1.2\n20,2.3\n30,3.4"))
	r := NewReader(f)
	var ints []int
	var floats []float32
	r.Loop(func(e struct {
		Int   int     `index:"0"`
		Float float32 `index:"1"`
	}) error {
		ints = append(ints, e.Int)
		floats = append(floats, e.Float)
		return nil
	})
	if err := r.Done(); err != nil {
		t.Error(err)
	}
	expectedInt := []int{10, 20, 30}
	expectedFloat := []float32{1.2, 2.3, 3.4}
	if !reflect.DeepEqual(expectedInt, ints) {
		t.Errorf("Unexpected %#v but got %#v", expectedInt, ints)
	}
	if !reflect.DeepEqual(expectedFloat, floats) {
		t.Errorf("Unexpected %#v but got %#v", expectedFloat, floats)
	}
}

func TestNewDecoder(t *testing.T) {
	d, err := newDecoder(reflect.TypeOf(struct {
		Name int `name:"name"`
		Age  int `name:"age"`
	}{}))
	if err != nil {
		t.Error(err)
	}
	if !d.needHeader() {
		t.Error("Unexpected")
	}
}

/*
func TestNewDecoder2(t *testing.T) {
	type mystruct struct {
		Name int `index:"1"`
		Age  int `index:"2"`
	}
	d, err := newDecoder(reflect.TypeOf(mystruct{}))
	if err != nil {
		t.Error(err)
	}
	if d.needHeader() {
		t.Error("Unexpected")
		return
	}
	var row mystruct
	err = d.decode([]string{"10", "30"}, reflect.ValueOf(&row))
	if err != nil {
		t.Error(err)
		return
	}
	t.Errorf("%#v", row)
}*/

/*
func Test3(t *testing.T) {
	f := bytes.NewReader([]byte("hello"))
	r := NewReader(f)
	var names []string
	r.Loop(func(e []string) error {
		for _, c := range e {
			names = append(names, c)
		}
		return nil
	})
	err := r.Done()
	if err != nil {
		t.Error(err)
	}
}

func Test4(t *testing.T) {
	f := bytes.NewReader([]byte("hello"))
	r := NewReader(f)
	var names []string
	r.Loop(func(e []string) error {
		for _, c := range e {
			names = append(names, c)
		}
		return nil
	})
	err := r.Done()
	if err != nil {
		t.Error(err)
	}
}

func Test1(t *testing.T) {
	f := bytes.NewReader([]byte("hello"))
	r := NewReader(f)
	var names []string
	var e struct {
		Name string
	}
	for r.Read(&e) {
		if e.Name == "" {
			break
		}
		names = append(names, e.Name)
	}
	err := r.Done()
	if err != nil {
		t.Error(err)
	}
}
*/
