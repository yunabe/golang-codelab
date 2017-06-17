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
	dec, err := newDecoder(in)
	if err != nil {
		r.err = err
		return
	}
	if dec.needHeader() {
		row, err := r.reader.Read()
		if err != nil {
			r.err = err
			return
		}
		err = dec.consumeHeader(row)
		if err != nil {
			r.err = err
			return
		}
	}
	for {
		row, err := r.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.err = err
			break
		}
		p := reflect.New(in)
		if err := dec.decode(row, p); err != nil {
			r.err = err
			break
		}
		rets := reflect.ValueOf(body).Call([]reflect.Value{p.Elem()})
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
		if cerr := r.closer.Close(); r.err == nil {
			r.err = cerr
		}
	}
	return r.err
}

type rowDecoder interface {
	decode(s []string, out reflect.Value) error
	needHeader() bool
	consumeHeader([]string) error
}

func createConverter(field reflect.StructField, enc string) (reflect.Value, error) {
	if field.Type.Kind() == reflect.Int {
		return reflect.ValueOf(strconv.Atoi), nil
	}
	if field.Type.Kind() == reflect.Float32 {
		return reflect.ValueOf(func(s string) (float32, error) {
			f, err := strconv.ParseFloat(s, 32)
			return float32(f), err
		}), nil
	}
	return reflect.ValueOf(nil), fmt.Errorf("Unexpected field type for %s: %s", field.Name, field.Type)

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
	conv, err := createConverter(field, enc)
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
	if t.Kind() != reflect.Struct {
		return nil, errors.New("error")
	}
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
