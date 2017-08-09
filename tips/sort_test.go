// This file demonstrate how to sort values in Go.

package tips

import (
	"sort"
	"testing"
)

type person struct {
	name string
	age  int
}

type byAge []person

func (a byAge) Len() int           { return len(a) }
func (a byAge) Less(i, j int) bool { return a[i].age < a[j].age }
func (a byAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func TestSortByAge(t *testing.T) {
	people := []person{
		{"Bob", 13},
		{"Alice", 8},
		{"Charlie", 20},
	}
	sort.Sort(byAge(people))
	if len(people) != 3 {
		t.Errorf("Unexpected length of people: %d", len(people))
		return
	}
	// Confirm the content of people is sorted.
	a := people[0]
	b := people[1]
	c := people[2]
	if a.name != "Alice" || a.age != 8 || b.name != "Bob" || b.age != 13 || c.name != "Charlie" || c.age != 20 {
		t.Errorf("people is not sorted correctly: %v", people)
	}
}
