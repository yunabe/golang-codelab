package easycsv

import (
	"bytes"
	"reflect"
	"testing"
)

func TestConverterInt(t *testing.T) {
	r := NewReader(bytes.NewBufferString("10,0xff,017"))
	var row []int
	ok := r.Read(&row)
	if !ok {
		t.Error("Read returned false unexpectedly")
	}
	expect := []int{10, 255, 15}
	if !reflect.DeepEqual(expect, row) {
		t.Errorf("Expected %v but got %v", expect, row)
	}
}

func TestConverterInvalid(t *testing.T) {
	r := NewReader(bytes.NewBufferString("hello"))
	var row []int
	ok := r.Read(&row)
	// TODO: Fix Reade to return false.
	if !ok {
		t.Error("Read returned false unexpectedly")
	}
}
