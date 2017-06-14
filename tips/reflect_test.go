package tips

// Reflection in Go:
// https://golang.org/pkg/reflect/

import (
	"reflect"
	"testing"
)

func TestReflectInt(t *testing.T) {
	i := 123
	ti := reflect.TypeOf(i)
	if ti.Kind() != reflect.Int {
		t.Errorf("Expected reflect.Int but got %v", ti.Kind())
		return
	}
	v := reflect.ValueOf(i)
	if v.Type() != ti {
		t.Errorf("Expected %v but got %v", ti, v.Type())
		return
	}
	// Retrieve the actual value.
	var ri int64 = v.Int()
	if ri != int64(i) {
		t.Errorf("Expected %d but got %d", i, ri)
	}
	// Interestingly, this does not panic :)
	if v.String() != "<int Value>" {
		t.Errorf("Expected an empty string but got %q", v.String())
	}
	// This panics! (using unaddressable value)
	// v.SetInt(456)
}

func TestReflectIntPointer(t *testing.T) {
	i := 123
	var p interface{} = &i
	ty := reflect.TypeOf(p)
	if ty.Kind() != reflect.Ptr {
		t.Errorf("Expected reflect.Ptr but got %v", ty.Kind())
		return
	}
	v := reflect.ValueOf(p)
	if v.Type() != ty {
		t.Errorf("Expected %v but got %v", ty, v.Type())
		return
	}
	var ev reflect.Value = v.Elem()
	if ev.Type().Kind() != reflect.Int {
		t.Errorf("Expected reflect.Int but got %v", ev.Type().Kind())
	}
	if ev.Kind() != reflect.Int {
		t.Errorf("Expected reflect.Int but got %v", ev.Kind())
	}
	if !ev.CanSet() {
		t.Error("CanSet must return true.")
	}
	// Set value!
	v.Elem().SetInt(456)
	if i != 456 {
		t.Errorf("Expected value: %d", i)
	}
}

func TestReflectStruct(t *testing.T) {
	type entry struct {
		name string
		Age  int `myattr:"myvalue"`
	}
	e := entry{name: "Alice", Age: 13}

	ty := reflect.TypeOf(e)
	if ty.Kind() != reflect.Struct {
		t.Errorf("Unexpected kind: %v", ty.Kind())
		return
	}
	if ty.NumField() != 2 {
		t.Errorf("Unexpected NumField(): %d", ty.NumField())
		return
	}
	var f0 reflect.StructField = ty.Field(0)
	if f0.Name != "name" || f0.Type.Kind() != reflect.String {
		t.Errorf("Unexpected field attrs: name = %q, type = %s", f0.Name, f0.Type)
	}
	f1 := ty.Field(1)
	if f1.Name != "Age" || f1.Type.Kind() != reflect.Int {
		t.Errorf("Unexpected field attrs: name = %q, type = %s", f1.Name, f1.Type)
	}
	// How to retrieve struct field tags.
	// https://golang.org/pkg/reflect/#StructTag
	if f1.Tag.Get("myattr") != "myvalue" {
		t.Errorf("Unexpected tag value: %q", f1.Tag.Get("myattr"))
	}

	// Value of struct.
	v := reflect.ValueOf(e)
	if v.Type() != ty {
		t.Errorf("Expected %v but got %v", ty, v.Type())
	}
	if v.NumField() != 2 {
		t.Errorf("Unexpected NumField(): %d", v.NumField())
		return
	}
	var vf0 reflect.Value = v.Field(0)
	if vf0.Type().Kind() != reflect.String {
		t.Errorf("Unexpected type for vf0: %v", vf0.Type())
	}
	if vf0.String() != "Alice" {
		t.Errorf("Unexpected value: %q", vf0.String())
	}
	vf1 := v.Field(1)
	if vf1.Type().Kind() != reflect.Int {
		t.Errorf("Unexpected type for vf1: %v", vf1.Type())
	}
	if vf0.String() != "Alice" {
		t.Errorf("Unexpected value: %q", vf0.String())
	}
	if vf1.Int() != 13 {
		t.Errorf("Unexpected value: %d", vf1.Int())
	}
	// fields of struct are not settable.
	if vf0.CanSet() || vf1.CanSet() {
		t.Error("CanSet of all fields must return false")
	}

	// This panics "using value obtained using unexported field".
	// vf0.SetString("hello")
	//
	// This panics with "using unaddressable value".
	// vf1.SetInt(20)
}

func TestReflectStructPointer(t *testing.T) {
	type entry struct {
		name string
		Age  int `myattr:"myvalue"`
	}
	p := &entry{name: "Alice", Age: 13}
	ty := reflect.TypeOf(p)
	if ty.Kind() != reflect.Ptr {
		t.Errorf("Unexpected kind: %v", ty.Kind())
		return
	}
	// Set field values.
	v := reflect.ValueOf(p)
	v.Elem().Field(1).SetInt(24)
	if p.Age != 24 {
		t.Errorf("Unexpected value: %v", p.Age)
	}
	// This panics "using value obtained using unexported field".
	// v.Elem().Field(0).SetString("Bob")
}

func TestReflectFunc(t *testing.T) {
	sum := func(x int, y int32) float32 {
		return float32(x + int(y))
	}
	var inter interface{} = sum
	ty := reflect.TypeOf(inter)
	if ty.Kind() != reflect.Func {
		t.Errorf("Unexpected kind: %v", ty.Kind())
		return
	}
	// Types of arguments and return values.
	if ty.NumIn() != 2 {
		t.Errorf("Unexpected NumIn() = %v", ty.NumIn())
	}
	if ty.NumOut() != 1 {
		t.Errorf("Unexpected NumOut() = %v", ty.NumOut())
	}
	if ty.In(0).Kind() != reflect.Int {
		t.Errorf("Unexpected first arg type: %v", ty.In(0).Kind())
	}
	if ty.Out(0).Kind() != reflect.Float32 {
		t.Errorf("Unexpected first arg type: %v", ty.Out(0).Kind())
	}

	// Invoke
	v := reflect.ValueOf(inter)
	returns := v.Call([]reflect.Value{
		reflect.ValueOf(12),
		reflect.ValueOf(int32(34)),
	})
	if len(returns) != 1 {
		t.Errorf("Unexpected size of slice was returned: %d", len(returns))
		return
	}
	f := returns[0].Float()
	if f != 46.0 {
		t.Errorf("Unexpected return value: %f", f)
	}

	// Call panics if arguments are invalid.
	// v.Call([]reflect.Value{})
}
