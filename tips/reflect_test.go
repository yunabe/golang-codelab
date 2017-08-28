package tips

// Reflection in Go:
// https://golang.org/pkg/reflect/

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"
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

func TestRelectEmbed(t *testing.T) {
	type embed struct {
		str string
		Str string
	}
	type Embed struct {
		Int     int
		integer int
	}
	type entry struct {
		F float32
		embed
		Embed
		float64
	}
	e := entry{
		F: 3.14,
		embed: embed{
			str: "hello",
			Str: "World",
		},
		Embed: Embed{
			Int:     100,
			integer: 10,
		},
		float64: 12.34,
	}
	ty := reflect.TypeOf(e)
	// Embeded struct are treated as fields as you may expect.
	if ty.NumField() != 4 {
		t.Errorf("Unexpected NumField() = %d", ty.NumField())
		return
	}
	f := ty.Field(0)
	if !(f.Name == "F" && !f.Anonymous) {
		t.Errorf("Unexpected: %#v", f)
	}
	f = ty.Field(1)
	// Anonymous is true if the field is an embedded field.
	if !(f.Name == "embed" && f.Anonymous) {
		t.Errorf("Unexpected: %#v", f)
	}
	f = ty.Field(2)
	// Anonymous is true if the field is an embedded field.
	if !(f.Name == "Embed" && f.Anonymous) {
		t.Errorf("Unexpected: %#v", f)
	}
	f = ty.Field(3)
	// Anonymous is true if the field is an embedded field.
	if !(f.Name == "float64" && f.Anonymous) {
		t.Errorf("Unexpected: %#v", f)
	}

	// TODO: Investigate Value.
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

func TestReflectSlice(t *testing.T) {
	slicePtr := reflect.New(reflect.SliceOf(reflect.TypeOf(int(0))))
	for i := 0; i < 10; i++ {
		slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(i*i)))
	}
	s := slicePtr.Elem().Interface().([]int)
	expected := []int{0, 1, 4, 9, 16, 25, 36, 49, 64, 81}
	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected %v but got %v", expected, s)
	}
}

func TestInterfaceToInterface(t *testing.T) {
	var err error
	var i interface{} = err
	if i != nil {
		t.Error("An interface from a nil interface must be nil too in Go.")
	}
	v := reflect.ValueOf(i)
	if v.IsValid() {
		t.Error("The value of nil must not be valid")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("v.IsNil() below must panic.")
			}
		}()
		// Wow, this panics!
		v.IsNil()
	}()
	func() {
		defer func() {
			p := recover()
			if p == nil {
				t.Error("The code must panic.")
				return
			}
			err := p.(error)
			if !strings.Contains(err.Error(), "interface is nil, not error") {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
		err = i.(error)
	}()
}

func TestNilAndCall(t *testing.T) {
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if errorType.Name() != "error" {
		t.Errorf("Unexpected errorType: %v", errorType)
	}

	fn := func(e error) error {
		if e != nil {
			panic("Call this function with nil!")
		}
		return nil
	}
	v := reflect.ValueOf(fn)

	// NG
	// arg := reflect.ValueOf(nil)
	arg := reflect.New(errorType).Elem()
	rets := v.Call([]reflect.Value{arg})
	if len(rets) != 1 {
		t.Errorf("Unexpected len(ret): %d", len(rets))
		return
	}
	ret := rets[0]
	// Unlike the test above, a nil interface in the return values is not invalid and has a type.
	if !ret.IsValid() {
		t.Error("ret is invalid unexpectedly")
	}
	if !ret.IsNil() {
		t.Error("ret.IsNil() returned false unexpectedly")
	}
	if ret.Type() != errorType {
		t.Errorf("Unexpected type: %v", ret.Type())
	}
}

type myInt int

func TestReflectNamedType(t *testing.T) {
	var a interface{} = myInt(3)
	va := reflect.ValueOf(a)
	ta := reflect.TypeOf(a)
	intType := reflect.TypeOf(0)
	if ta.String() != "tips.myInt" {
		t.Errorf("ta.String() must be \"tips.myInt\", but got %q", ta.String())
	}
	if ta == intType {
		t.Errorf("typeof(myInt) == typeof(int) must be false")
	}
	if ta.Name() != "myInt" {
		t.Errorf("ta.Name() must be \"myInt\", but got %q", ta.Name())
	}
	if va.Kind() != reflect.Int {
		// Although typeof(myInt) != typeof(int), Kind() returns reflect.Int.
		t.Error("va.Kind() must be reflect.Int")
	}
	if va.Int() != 3 {
		t.Errorf("va.Int() must be 3, but got %v", va.Int())
	}
}

func TestReflectAddr(t *testing.T) {
	v := reflect.ValueOf(10)
	if v.CanAddr() {
		t.Errorf("v.CanAddr() must return false")
	}
	v = reflect.New(reflect.TypeOf(0))
	if !v.Elem().CanAddr() {
		t.Errorf("v.Elem().CanAddr() must return true")
	}
}

type childField struct {
	s string
}

type structWithUnexported struct {
	i     int8
	f     float32
	child childField
}

func expectUnexportedError(t *testing.T) {
	r := recover()
	if r == nil {
		t.Errorf("Expected panic but no panic was thrown")
	}
	if msg, _ := r.(string); !strings.Contains(msg, "unexported field") {
		t.Errorf("Unexpected error: %v", r)
	}
}

func TestStructUnexportedField(t *testing.T) {
	u := structWithUnexported{12, 3.4, childField{"hello"}}
	v := reflect.ValueOf(&u).Elem()

	// You can read unexported fields with reflect.
	if v.Field(0).Kind() != reflect.Int8 {
		t.Errorf("Expected %v but got %v", reflect.Int8, v.Field(0).Kind())
	}
	i := int8(v.Field(0).Int())
	if i != u.i {
		t.Errorf("Expected %v but got %v", u.i, i)
	}
	s := v.FieldByName("child").FieldByName("s").String()
	if s != u.child.s {
		t.Errorf("Expected %q but got %q", u.child.s, s)
	}

	// But you can not use Interface() to read unexported fields for some reasons.
	func() {
		defer expectUnexportedError(t)
		v.FieldByName("i").Interface()
	}()

	// You can not set unexported fields normally.
	func() {
		defer expectUnexportedError(t)
		v.FieldByName("f").SetFloat(5.6)
	}()
	// But, you can set values to unexported fields with unsafe.
	// Ref: https://stackoverflow.com/questions/42664837/access-unexported-fields-in-golang-reflect
	ff := v.FieldByName("f")
	newF := float32(7.8)
	reflect.NewAt(ff.Type(), unsafe.Pointer(ff.UnsafeAddr())).Elem().SetFloat(float64(newF))
	if u.f != newF {
		t.Errorf("Expected %v but got %v", newF, u.f)
	}
}
