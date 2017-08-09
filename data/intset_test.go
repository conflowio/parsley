// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntSetShouldCreateAnEmptySet(t *testing.T) {
	i := NewIntSet()
	assert.Equal(t, []int{}, i.data)
}

func TestNewIntSetShouldInsertValues(t *testing.T) {
	i := NewIntSet(3, 2, 2, 1, 0)
	assert.Equal(t, []int{0, 1, 2, 3}, i.data)
}

func TestEmptyIntSetReturnsWithAnEmptyIntSet(t *testing.T) {
	i := EmptyIntSet()
	assert.Equal(t, []int{}, i.data)
}

func TestLenShouldReturnWithSetLength(t *testing.T) {
	assert.Equal(t, 0, IntSet{[]int{}}.Len())
	assert.Equal(t, 1, IntSet{[]int{1}}.Len())
	assert.Equal(t, 2, IntSet{[]int{1, 2}}.Len())
}

func TestInsert(t *testing.T) {
	type TC struct {
		name string
		i1   []int
		val  int
		i2   []int
	}
	testCases := []TC{
		{"Insert to empty list", []int{}, 1, []int{1}},
		{"Insert zero", []int{}, 0, []int{0}},
		{"Existing item should not be duplicated", []int{1}, 1, []int{1}},
		{"Insert to the end of the list", []int{1}, 2, []int{1, 2}},
		{"Insert to the beginning of the list", []int{2}, 1, []int{1, 2}},
		{"Insert to the middle of the list", []int{1, 3}, 2, []int{1, 2, 3}},
	}
	for _, tc := range testCases {
		i1c := make([]int, len(tc.i1))
		copy(i1c, tc.i1)
		i2 := IntSet{tc.i1}.Insert(tc.val)
		assert.Equal(t, i1c, tc.i1, fmt.Sprintf("Set is not immutable: %s", tc.name))
		assert.Equal(t, IntSet{tc.i2}, i2, fmt.Sprintf("Failed: %s", tc.name))
	}
}

func TestUnion(t *testing.T) {
	type TC struct {
		name string
		i1   []int
		i2   []int
		i3   []int
	}
	testCases := []TC{
		{"Empty lists", []int{}, []int{}, []int{}},
		{"First empty", []int{}, []int{1}, []int{1}},
		{"Second empty", []int{1}, []int{}, []int{1}},
		{"First longer", []int{1, 3}, []int{2}, []int{1, 2, 3}},
		{"Second longer", []int{2}, []int{1, 3}, []int{1, 2, 3}},
		{"First then second", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"Second then first", []int{3, 4}, []int{1, 2}, []int{1, 2, 3, 4}},
		{"Mixed", []int{1, 3}, []int{2, 4}, []int{1, 2, 3, 4}},
		{"Ignore duplicated", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
	}
	for _, tc := range testCases {
		i1c := make([]int, len(tc.i1))
		copy(i1c, tc.i1)
		i2c := make([]int, len(tc.i2))
		copy(i2c, tc.i2)
		i3 := IntSet{tc.i1}.Union(IntSet{tc.i2})
		assert.Equal(t, i1c, tc.i1, fmt.Sprintf("Set is not immutable: %s", tc.name))
		assert.Equal(t, i2c, tc.i2, fmt.Sprintf("Set is not immutable: %s", tc.name))
		assert.Equal(t, IntSet{tc.i3}, i3, fmt.Sprintf("Failed: %s", tc.name))
	}
}

func TestEachShouldNotCallFunctionForEmptyList(t *testing.T) {
	calledValue := -1
	f := func(val int) {
		calledValue = val
	}
	i := NewIntSet()
	i.Each(f)
	assert.Equal(t, -1, calledValue)
}

func TestEachShouldCallFunctionForAllItems(t *testing.T) {
	calledValues := []int{}
	f := func(val int) {
		calledValues = append(calledValues, val)
	}
	i := IntSet{[]int{1, 2}}
	i.Each(f)
	assert.Equal(t, []int{1, 2}, calledValues)
}
