package tips

import (
	"testing"
)

type interfaceX interface {
	implementX()
}

type interfaceY interface {
	implementY()
}

type interfaceZ interface {
	interfaceX
	interfaceY
}

type impl int

func (impl) implementX() {}
func (impl) implementY() {}

func TestInterfaceEqualityOp(t *testing.T) {
	// References:
	// https: //golang.org/ref/spec#Comparison_operators

	// 1. Two interface values are equal if they have identical dynamic types and equal dynamic values or if both have value nil.
	var a, b, c interface{}
	a = 10
	b = 20
	c = 10
	if a == b {
		t.Errorf("a != b must be true")
	}
	if a != c {
		t.Errorf("a == c must be true")
	}

	var d interfaceX = impl(10)
	var f interfaceY = impl(10)
	var g interfaceZ = impl(10)
	var h = impl(10)

	// compile error: invalid operator
	// if d != f {}
	if interface{}(d) != interface{}(f) {
		t.Errorf("d == f must be true")
	}

	if d != g {
		t.Errorf("d == g must be true")
	}
	if f != g {
		t.Errorf("d == g must be true")
	}
	// 2. A value x of non-interface type X and a value t of interface type T are comparable when values of type X are comparable and X implements T.
	if h != d || h != f || h != g {
		t.Errorf("h must equals to d, f and g")
	}

	// 3. A comparison of two interface values with identical dynamic types causes a run-time panic if values of that type are not comparable.
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("The following code block must panic")
			}
		}()
		var s0, s1 interface{}
		s0 = []int{1, 2, 3}
		s1 = []int{1, 2, 3}
		t.Errorf("s0 == s1 must panic but got %v", s0 == s1)
	}()
}
