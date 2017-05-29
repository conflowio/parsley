package data

import (
	"sort"
)

// IntSet is a sorted immutable integer set
type IntSet struct {
	data []int
}

// NewIntSet creates a new integer set
func NewIntSet(values ...int) IntSet {
	i := IntSet{make([]int, 0, len(values))}
	for _, val := range values {
		i.insertValue(val)
	}
	return i
}

// Len returns with the length of the set
func (i IntSet) Len() int {
	return len(i.data)
}

// Insert adds a new item to the set
func (i IntSet) Insert(val int) IntSet {
	if len(i.data) == 0 {
		return IntSet{[]int{val}}
	}
	i2 := i
	i2.insertValue(val)
	return i2
}

func (i *IntSet) insertValue(val int) {
	index := sort.SearchInts(i.data, val)
	if index < len(i.data) && i.data[index] == val {
		return
	}
	i.data = append(i.data, 0)
	copy(i.data[index+1:], i.data[index:])
	i.data[index] = val
}

// Union returns with the union of the two set
func (i IntSet) Union(i2 IntSet) IntSet {
	if len(i2.data) == 0 {
		return i
	} else if len(i.data) == 0 {
		return i2
	}
	i3 := IntSet{make([]int, 0, len(i.data)+len(i2.data))}
	var n1, n2 int
	for n1 < len(i.data) || n2 < len(i2.data) {
		if n2 >= len(i2.data) || n1 < len(i.data) && i.data[n1] < i2.data[n2] {
			i3.data = append(i3.data, i.data[n1])
			n1++
		} else if n1 >= len(i.data) || n2 < len(i2.data) && i2.data[n2] < i.data[n1] {
			i3.data = append(i3.data, i2.data[n2])
			n2++
		} else {
			i3.data = append(i3.data, i.data[n1])
			n1++
			n2++
		}
	}
	return i3
}

// Each runs the given function on all elements of the set
func (i IntSet) Each(f func(val int)) {
	for _, v := range i.data {
		f(v)
	}
}
