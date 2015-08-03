package main

import (
	"fmt"
	"sort"
)

type person struct {
	name string
	age int
}

type byAge []person

func (a byAge)  Len() int { return len(a) }
func (a byAge)  Less(i, j int) bool { return a[i].age < a[j].age }
func (a byAge)  Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func main() {
	people := []person{
		{"Bob", 13},
		{"Alice", 8},
		{"Charlie", 20},
	}
	sort.Sort(byAge(people))
	for i, p := range people {
		fmt.Printf("index: %d, name: %s, age: %d\n", i, p.name, p.age)
	}
}
