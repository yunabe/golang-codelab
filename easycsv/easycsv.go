package easycsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var Break = errors.New("break")

type Reader struct {
	reader *csv.Reader
	closer io.Closer
	done   bool
	err    error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: csv.NewReader(r),
	}
}

func NewReadCloser(r io.ReadCloser) *Reader {
	return &Reader{
		reader: csv.NewReader(r),
		closer: r,
	}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (r *Reader) Loop(body interface{}) {
	if r.err != nil {
		return
	}
	if body == nil {
		r.err = errors.New("The argument of Loop must not be nil.")
		return
	}
	v := reflect.TypeOf(body)
	if v.Kind() != reflect.Func {
		r.err = fmt.Errorf("The argument of Loop must be func but got %v", v.Kind())
		return
	}
	if v.NumIn() != 1 || v.NumOut() != 1 {
		r.err = fmt.Errorf("The function passed to Loop must return receive one argument and return one argument")
		return
	}
	out := v.Out(0)
	if out.Kind() != reflect.Bool && out != errorType {
		r.err = fmt.Errorf("The function passed to Loop must return error or bool")
		return
	}
	in := v.In(0)
	if in.Kind() != reflect.Struct {
		r.err = fmt.Errorf("The function passed to Loop must receive a struct")
		return
	}
	numf := in.NumField()
	if numf == 0 {
		r.err = errors.New("The struct passed to Loop must have at least one field")
		return
	}
	p := reflect.New(in)
	var privates []string
	for i := 0; i < p.Elem().NumField(); i++ {
		if !p.Elem().Field(i).CanSet() {
			privates = append(privates, in.Field(i).Name)
		}
	}
	if privates != nil {
		r.err = fmt.Errorf("The struct passed to Loop must not have private fields: %s", strings.Join(privates, ", "))
		return
	}
}

func (r *Reader) Read(e interface{}) bool {
	return true
}

func (r *Reader) Done() error {
	if r.done {
		return r.err
	}
	r.done = true
	if r.closer != nil {
		if cerr := r.closer.Close(); r.err != nil {
			r.err = cerr
		}
	}
	return r.err
}
