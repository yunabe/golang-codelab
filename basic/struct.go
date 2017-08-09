package main

import (
	"fmt"
)

type myStruct struct {
	val int
}

func (s myStruct) Method0() {
}

func (s *myStruct) Method1() {}

type myStruct2 myStruct
type myStruct3 *myStruct
type myStruct4 struct {
	myStruct
}
type myStruct5 struct {
	*myStruct
}

func playWithMyStruct() {
	s2 := myStruct2{}
	fmt.Println(s2.val)
	// s2.Method0()
	// s2.Method1()

	s3 := myStruct3(nil)
	fmt.Println(s3.val)
	// s2.Method0()
	// s2.Method1()

	s4 := myStruct4{
		myStruct: myStruct{},
	}
	s4.Method1()
	s4.Method0()

	s5 := myStruct5{
		myStruct: &myStruct{},
	}
	s5.Method1()
	s5.Method0()
}
