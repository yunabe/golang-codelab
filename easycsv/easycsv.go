package easycsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var Break = errors.New("break")

type Reader struct {
	// csv.Reader. To read content from csv, use readLine.
	csv    *csv.Reader
	closer io.Closer
	done   bool
	// An error occurred while processing csv. io.EOF is stored when csv is reached to the end.
	err error

	// Used from readLine.
	lineno    int
	firstLine []string
	cur       []string
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		csv: csv.NewReader(r),
	}
}

func NewReadCloser(r io.ReadCloser) *Reader {
	return &Reader{
		csv:    csv.NewReader(r),
		closer: r,
	}
}

// readLine reads a line from r.csv and update r.err, r.cur, r.lineno and r.firstLine.
// readLine does not update r.err. io.EOF is returned when csv reached to the end.
func (r *Reader) readLine() {
	line, err := r.csv.Read()
	if err != nil {
		r.err = err
		return
	}
	r.cur = line
	r.lineno++
	if r.lineno == 1 {
		r.firstLine = line
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
	var inStruct reflect.Type
	if in.Kind() == reflect.Struct {
		inStruct = in
	} else if in.Kind() == reflect.Ptr && in.Elem().Kind() == reflect.Struct {
		inStruct = in.Elem()
	} else if in.Kind() == reflect.Slice {
		inStruct = in
	} else {
		r.err = fmt.Errorf("The function passed to Loop must receive a struct, a pointer to a struct or a slice")
		return
	}
	if in.Kind() != reflect.Slice {
		numf := inStruct.NumField()
		if numf == 0 {
			r.err = errors.New("The struct passed to Loop must have at least one field")
			return
		}
	}
	dec, err := newDecoder(inStruct)
	if err != nil {
		r.err = err
		return
	}
	if dec.needHeader() {
		if r.lineno == 0 {
			// Loop quits immediately if the csv is empty.
			r.readLine()
			if r.err != nil {
				return
			}
		}
		err = dec.consumeHeader(r.firstLine)
		if err != nil {
			r.err = err
			return
		}
	}
	for {
		r.readLine()
		if r.err != nil {
			break
		}
		p := reflect.New(inStruct)
		if err := dec.decode(r.cur, p); err != nil {
			r.err = err
			break
		}
		arg := p
		if in.Kind() == reflect.Struct || in.Kind() == reflect.Slice {
			arg = p.Elem()
		}
		rets := reflect.ValueOf(body).Call([]reflect.Value{arg})
		if len(rets) == 0 || rets[0].IsNil() {
			continue
		}
		if rets[0].Kind() == reflect.Bool {
			cont := rets[0].Bool()
			if cont {
				continue
			} else {
				break
			}
		}
		err = rets[0].Interface().(error)
		if err == Break {
			break
		}
		r.err = err
		// TODO: break here to terminate the loop on error.
	}
}

func (r *Reader) Read(e interface{}) bool {
	if r.err != nil {
		return false
	}
	t := reflect.TypeOf(e)
	if t.Kind() != reflect.Ptr {
		r.err = fmt.Errorf("The argument of Read must be a pointer to a struct or a pointer to a slice, but got %v", t.Kind())
		return false
	}
	if t.Elem().Kind() != reflect.Struct && t.Elem().Kind() != reflect.Slice {
		r.err = fmt.Errorf("The argument of Read must be a pointer to a struct or a pointer to a slice, but got a pointer to %v", t.Elem().Kind())
		return false
	}
	decoder, err := newDecoder(t.Elem())
	if err != nil {
		r.err = err
		return false
	}
	if decoder.needHeader() {
		if r.lineno == 0 {
			// Loop quits immediately if the csv is empty.
			r.readLine()
			if r.err != nil {
				return false
			}
		}
		decoder.consumeHeader(r.firstLine)
	}
	r.readLine()
	if r.err != nil {
		return false
	}
	// TODO: Reset with zero.
	decoder.decode(r.cur, reflect.ValueOf(e))
	return true
}

func (r *Reader) nonEOFError() error {
	if r.err == nil || r.err == io.EOF {
		return nil
	}
	return r.err
}

// Done returns the first non-EOF error that was encountered by the Reader.
// Done also closes the internal Closer if the Reader is instantiated with NewReaderCloser.
func (r *Reader) Done() error {
	if r.done {
		return r.nonEOFError()
	}
	r.done = true
	if r.closer != nil {
		if cerr := r.closer.Close(); r.err == nil || r.err == io.EOF {
			r.err = cerr
		}
	}
	return r.nonEOFError()
}

// DoneDefer do the same thing as Done does. But it outputs an error to the argument.
func (r *Reader) DoneDefer(err *error) {
	e := r.Done()
	if *err == nil && e != nil {
		*err = e
	}
}

type rowDecoder interface {
	decode(s []string, out reflect.Value) error
	needHeader() bool
	consumeHeader([]string) error
}

func parseStructTag(field reflect.StructField,
	fieldIdx int,
	nameMap map[string]int,
	idxMap map[int]int,
	converters *[]reflect.Value,
	errors *[]string) {
	tag := field.Tag
	name := tag.Get("name")
	index := tag.Get("index")
	if name == "" && index == "" {
		*errors = append(*errors, fmt.Sprintf("Please specify name or index to the struct field: %s", field.Name))
		return
	}
	if name != "" && index != "" {
		*errors = append(*errors, fmt.Sprintf("Please specify name or index to the struct field: %s", field.Name))
		return
	}
	enc := tag.Get("enc")
	conv, err := createDefaultConverter(field, enc)
	if err != nil {
		*errors = append(*errors, err.Error())
		return
	}
	*converters = append(*converters, conv)
	if name != "" {
		nameMap[name] = fieldIdx
		return
	}
	i, err := strconv.Atoi(index)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("Failed to parse index of field %s: %q", field.Name, index))
		return
	}
	idxMap[i] = fieldIdx
}

func newDecoder(t reflect.Type) (rowDecoder, error) {
	if t.Kind() == reflect.Struct {
		return newStructDecoder(t)
	} else if t.Kind() == reflect.Slice {
		return newSliceDecoder(t)
	}
	panic("newDecoder must be called with struct or slice.")
}

func newSliceDecoder(t reflect.Type) (rowDecoder, error) {
	elem := t.Elem()
	c := createDefaultConverterFromType(elem)
	if !c.IsValid() {
		return nil, fmt.Errorf("Failed to create a converter for %v", t)
	}
	return &sliceRowDecoder{
		elemType:  elem,
		converter: c,
	}, nil
}

type sliceRowDecoder struct {
	elemType  reflect.Type
	converter reflect.Value
}

func (d *sliceRowDecoder) needHeader() bool             { return false }
func (d *sliceRowDecoder) consumeHeader([]string) error { return nil }
func (d *sliceRowDecoder) decode(s []string, out reflect.Value) error {
	slicePtr := reflect.New(reflect.SliceOf(d.elemType))
	for _, e := range s {
		rets := d.converter.Call([]reflect.Value{reflect.ValueOf(e)})
		// TODO: Handle error.
		slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), rets[0]))
	}
	out.Elem().Set(slicePtr.Elem())
	return nil
}

func newStructDecoder(t reflect.Type) (rowDecoder, error) {
	if t.NumField() == 0 {
		return nil, errors.New("The struct has no field")
	}
	v := reflect.New(t).Elem()
	var unexported []string
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			unexported = append(unexported, t.Field(i).Name)
		}
	}
	if unexported != nil {
		return nil, fmt.Errorf("The struct passed to Loop must not have unexported fields: %s", strings.Join(unexported, ", "))
	}

	var tagErrors []string
	nameMap := make(map[string]int)
	idxMap := make(map[int]int)
	var converters []reflect.Value
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		parseStructTag(f, i, nameMap, idxMap, &converters, &tagErrors)
	}
	if len(nameMap) != 0 && len(idxMap) != 0 {
		tagErrors = append(tagErrors, "Fields with name and fields with index are mixed")
	}
	if tagErrors != nil {
		return nil, errors.New(strings.Join(tagErrors, "\n"))
	}
	if len(converters) != t.NumField() {
		panic("converters size mismatch")
	}
	if len(nameMap) != 0 {
		idxMap = nil
	} else {
		nameMap = nil
	}
	return &structRowDecoder{
		structType: t,
		converters: converters,
		names:      nameMap,
		indice:     idxMap,
	}, nil
}

type structRowDecoder struct {
	structType reflect.Type
	converters []reflect.Value
	names      map[string]int
	indice     map[int]int
}

func (d *structRowDecoder) consumeHeader(header []string) error {
	indice := make(map[int]int)
	for i, col := range header {
		idx, ok := d.names[col]
		if !ok {
			continue
		}
		indice[i] = idx
		delete(d.names, col)
	}
	d.indice = indice
	if len(d.names) != 0 {
		var unused []string
		for n := range d.names {
			unused = append(unused, n)
		}
		return fmt.Errorf("%s did not appear in the first line", strings.Join(unused, ", "))
	}
	d.names = nil
	return nil
}

func (d *structRowDecoder) decode(row []string, out reflect.Value) error {
	for i, j := range d.indice {
		rets := d.converters[j].Call([]reflect.Value{reflect.ValueOf(row[i])})
		if len(rets) != 2 {
			panic("converter must return two values.")
		}
		if !rets[1].IsNil() {
			return rets[1].Interface().(error)
		}
		out.Elem().Field(j).Set(rets[0])
	}
	return nil
}

func (d *structRowDecoder) needHeader() bool {
	return d.names != nil
}
