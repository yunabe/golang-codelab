package main

import (
	"fmt"
)

func fact(n int) int {
  if n > 1 {
    return n * fact(n - 1)
  }
  return 1
}

func play_with_array() {
	fmt.Println("### play_with_array ###")
	var a [3]int  // all set to 0
	fmt.Println(a)
	fmt.Printf("a[0] == %d\n", a[0])
	fmt.Printf("len(a) == %d\n", len(a))

	// Array literals
	a  = [3]int{1, 2, 3}
	fmt.Println(a)
	a  = [...]int{4, 5, 6}
	fmt.Println(a)
	a  = [3]int{0:7, 2:8}
	fmt.Println(a)

	// Type error:
	// cannot use array literal (type [4]int)
	// as type [3]int in assignment
	// a = [4]int{1, 2, 3, 4}
}

func play_with_slice() {
	fmt.Println("### play_with_slice ###")
	var array [10]int;
	array[0] = 0
	array[1] = 1
	array[2] = 2

	// slice
	var slice []int
	fmt.Println("initial value of slice <")
	fmt.Print("  slice == ")
	fmt.Println(slice)  // an empty slice
	// an empty slice (an initial value of slice) is nil.
	// slice is a kind of pointer.
	fmt.Printf("  (slice == nil) == %t\n", slice == nil)  // true
	fmt.Printf("  len(slice) == %d\n", len(slice))  // 0
	fmt.Printf("  cap(slice) == %d\n", cap(slice))  // 0
	fmt.Println(">")
	
	slice = array[:]
	fmt.Printf("len(slice) = %d\n", len(slice))
	fmt.Printf("cap(slice) = %d\n", cap(slice))
	fmt.Printf("slice[0] = %d\n", slice[0])
	slice = array[1:]
	fmt.Printf("len(slice) = %d\n", len(slice))
	fmt.Printf("cap(slice) = %d\n", cap(slice))
	fmt.Printf("slice[0] = %d\n", slice[0])
	
	// slice literal
	slice = []int{1, 3, 5}
	fmt.Println(slice)

	// alloc arbitrary size array to slice
	slice = make([]int, 20)
	fmt.Println("** slice = make([]int, 20) **")
	fmt.Printf("len(slice) == %d\n", len(slice))
	fmt.Printf("cap(slice) == %d\n", cap(slice))

	slice = array[3:5]
	fmt.Printf("len(slice) == %d\n", len(slice))
	fmt.Printf("cap(slice) == %d\n", cap(slice))

	slice = make([]int, 10, 20)
	fmt.Printf("len(slice) == %d\n", len(slice))
	fmt.Printf("cap(slice) == %d\n", cap(slice))
	// fmt.Println(slice[15])  // index out of range
	slice = slice[:20]
	fmt.Printf("slice[15] = %d\n", slice[15])
	// slice = slice[:100]  // index out of range
}

func play_with_append() {
	// There is a special function "append" in golang.
	// "apend" is special as types of return values are changed by args of
	// "append". We can not define such kind of functions in golang.
	//
	//     C++        :     Go
	// --------------------------------
	// vector<T> v    : var v T[]
	// v.push_back(e) : v = append(v, e)
	// v.size()       : len(v)
	// v[i]           : v[i]
	// vector<T>*  p  : var p *T[]

	fmt.Println("### play_with_append ###")
	ia := []int{1, 2, 3}
	fmt.Println("ia =", ia)
	ia = append(ia, 4, 5, 6)
	fmt.Println("ia =", ia)

	sa := []string{"foo", "bar"}
	fmt.Println("sa =", sa)
	sa = append(sa, "hoge", "piyo")
	fmt.Println("sa =", sa)
}

func play_with_pointer() {
	fmt.Println("### play_with_pointer ###")
	var x int = 1
	var p *int = &x
	fmt.Printf("*p == %d\n", *p)
	x = 2
	fmt.Printf("*p == %d\n", *p)
	p = nil
	fmt.Println(p)
	fmt.Println(p == nil)
}

func play_with_map() {
  fmt.Println("### play_with_map ###")
	// Equivalent to m := make(map[string]float64)
	var m map[string]float64 = make(map[string]float64)
	m["foo"] = 10.0
	var m2 = m
	m2["foo"] = 20.0
	fmt.Println(m["foo"])  // Shows 20.0. map is a thin object like slice or pointer.

	// nil is an default initial value.
	// map is also a kind of pointer like slice.
	m2 = nil
	fmt.Println(m2)
	fmt.Println(m2 == nil)

	m3 := map[string]float64{"foo": 1, "bar": 3}  // map literal
	fmt.Println(m3)
	var value = m3["piyo"]  // No exception. value is set to 'zero'.
	fmt.Printf("value == %f\n", value)
	value, ok := m3["piyo"]  // ok is set true if a map has a key.
	fmt.Println("value, ok =", value, ok)

	delete(m3, "foo")  // Removes an element.
	fmt.Println(m3)
}

func play_with_if() {
	fmt.Println("### play_with_if ###")
	n := 2
	m := 3
	if n % 2 == 1 {
		fmt.Println("cond0")
	} else if m == 3 {
		fmt.Println("cond1")
	} else {
		fmt.Println("cond2")
	}

	// multiline condition
	if n == 2 &&
		m == 3 {
		// n == 2
		// && m == 3 -> syntax error: unexpected &&, expecting {
	}

	// We can use 1 initialization sentence with ;
	// Of course, vars defined here are visible only from inside of if-sentence.
	if tmp := n * m; tmp == 6 {
		// var tmp = ... is not valid.
	}
}

func play_with_for() {
	fmt.Println("### play_with_for ###")
	for i := 0; i < 3; i++ {
		fmt.Printf("usual for: i == %d\n", i)
	}
	var count int
	count = 0
	for count < 3 {
		fmt.Printf("while loop: count == %d\n", count)
		count++
	}
	count = 0
	for {
		fmt.Printf("infinite loop: count == %d\n", count)
		count++
		if count > 3 {
			break
		}
	}
}

func play_with_foreach() {
	fmt.Println("### play_with_foreach ###")
	array := [...]string{"foo", "bar", "baz"}
	for i, v := range array {
		fmt.Printf("for range with array: %d - %s\n", i, v)
	}
	dict := map[string]int{"foo": 3, "bar": 4}
	for key, value := range dict {
		fmt.Printf("for range with map: %s - %d\n", key, value)
	}
}

func play_with_switch() {
	fmt.Println("### play_with_switch ###")
	for i := 0; i < 4; i++ {
		fmt.Println("** i ==", i, "**")
		switch i {
		case 0:
			fmt.Println("i is 0.")
			// fall through is not default behavior of switch in golang.
			fallthrough
		case 1, 2:
			fmt.Println("i is 0, 1 or 2.")
		default:
			fmt.Println("i is not either 0, 1 or 2.")
		}
	}
	// In golang, we can switch instead of if-else.
	for i := -1; i <= 1; i++ {
		fmt.Println("** i ==", i, "**")
		switch {
		case i < 0:
			fmt.Println("i is negative.")
		case i == 0: // Do nothing
		default:
			fmt.Println("i is positive.")
		}
	}

	// With initialization
	switch value := 10; value % 2 {
	case 0:
		fmt.Println("value is even.")
	case 1:
		fmt.Println("value is odd.")
	}
	// Hmm..., I feel it is slightly confusing...
	switch value := 7; {
	case value % 2 == 0:
		fmt.Println("value is even.")
	default:
		fmt.Println("value is odd.")
	}
	// Also, we can use "type switch" in golang. See ooo.go for details.
}

func play_with_label() {
	fmt.Println("play_with_label")
	OuterLoop: for i := 0;; i++ {
		for j := 0; j <= i; j++ {
			fmt.Println("i =", i, ", j =", j)
			if i * j == 6 {
				fmt.Println("break!")
				  break OuterLoop
			}
		}
  }
}

func basicMain() {
	n := 5
	fmt.Printf("fact(%d) = %d\n", n, fact(n))
	play_with_array()
	play_with_slice()
	play_with_append()
	play_with_pointer()
	play_with_map()
	play_with_if()
	play_with_for()
	play_with_foreach()
	play_with_switch()
	play_with_label()
}
