package easycsv

import (
	"strings"
	"testing"
)
import "bytes"

func TestLoopNil(t *testing.T) {
	f := bytes.NewReader([]byte(""))
	r := NewReader(f)
	r.Loop(nil)
	err := r.Done()
	if err == nil || strings.Contains("adafdsa", err.Error()) {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLoop(t *testing.T) {
	f := bytes.NewReader([]byte(""))
	r := NewReader(f)
	var names []string
	r.Loop(func(e struct {
		Name string `name:"name"`
		Age  int    `name:"age",enc:"hex"`
	}) error {
		if e.Name == "" {
			return Break
		}
		names = append(names, e.Name)
		return nil
	})
	err := r.Done()
	if err != nil {
		t.Error(err)
	}
}

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
