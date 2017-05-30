package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntMapShouldCreateAnEmptyMap(t *testing.T) {
	i := NewIntMap(nil)
	assert.Equal(t, map[int]int{}, i.data)
}

func TestNewIntMapShouldCreateMapWithGivenData(t *testing.T) {
	data := map[int]int{1: 2}
	i := NewIntMap(data)
	assert.Equal(t, data, i.data)
}

func TestInc(t *testing.T) {
	type TC struct {
		name string
		i1   map[int]int
		val  int
		i2   map[int]int
	}
	testCases := []TC{
		TC{"Insert to empty map", map[int]int{}, 1, map[int]int{1: 1}},
		TC{"Insert new value", map[int]int{1: 1}, 2, map[int]int{1: 1, 2: 1}},
		TC{"Increase existing value", map[int]int{1: 1, 2: 1}, 2, map[int]int{1: 1, 2: 2}},
	}
	for _, tc := range testCases {
		i1c := make(map[int]int, len(tc.i1))
		for k, v := range tc.i1 {
			i1c[k] = v
		}
		i2 := IntMap{tc.i1}.Inc(tc.val)
		assert.Equal(t, i1c, tc.i1, fmt.Sprintf("Map is not immutable: %s", tc.name))
		assertMapsAreEqual(t, IntMap{tc.i2}, i2, fmt.Sprintf("Failed: %s", tc.name))
	}
}

func TestFilter(t *testing.T) {
	type TC struct {
		name string
		i1   map[int]int
		keys []int
		i2   map[int]int
	}
	testCases := []TC{
		TC{"Empty keys", map[int]int{1: 2}, []int{}, map[int]int{}},
		TC{"Empty map", map[int]int{}, []int{1}, map[int]int{}},
		TC{"None of the keys present", map[int]int{1: 2}, []int{2}, map[int]int{}},
		TC{"Only include keys", map[int]int{1: 2, 3: 4, 5: 6}, []int{3, 5}, map[int]int{3: 4, 5: 6}},
	}
	for _, tc := range testCases {
		i1c := make(map[int]int, len(tc.i1))
		for k, v := range tc.i1 {
			i1c[k] = v
		}
		i2 := IntMap{tc.i1}.Filter(IntSet{tc.keys})
		assert.Equal(t, i1c, tc.i1, fmt.Sprintf("Map is not immutable: %s", tc.name))
		assertMapsAreEqual(t, IntMap{tc.i2}, i2, fmt.Sprintf("Failed: %s", tc.name))
	}
}

func TestGetReturnsZeroForNonExistingItem(t *testing.T) {
	i := NewIntMap(nil)
	assert.Exactly(t, 0, i.Get(1))
}

func TestGetReturnsValue(t *testing.T) {
	i := IntMap{map[int]int{1: 2, 3: 4}}
	assert.Exactly(t, 4, i.Get(3))
}

func TestEachShouldNotCallFunctionForEmptyMap(t *testing.T) {
	called := false
	f := func(key int, val int) {
		called = true
	}
	i := NewIntMap(nil)
	i.Each(f)
	assert.False(t, called)
}

func TestEachShouldCallFunctionForAllMapItems(t *testing.T) {
	calledKeys := []int{}
	calledValues := []int{}
	f := func(key int, val int) {
		calledKeys = append(calledKeys, key)
		calledValues = append(calledValues, val)
	}
	i := IntMap{map[int]int{1: 2, 3: 4}}
	i.Each(f)
	assertSlicesAreEqual(t, []int{1, 3}, calledKeys)
	assertSlicesAreEqual(t, []int{2, 4}, calledValues)
}

func TestKeys(t *testing.T) {
	type TC struct {
		name string
		i    map[int]int
		keys []int
	}
	testCases := []TC{
		TC{"Empty map", map[int]int{}, []int{}},
		TC{"Non-empty map", map[int]int{1: 2, 3: 4}, []int{1, 3}},
	}
	for _, tc := range testCases {
		assertSlicesAreEqual(t, tc.keys, IntMap{tc.i}.Keys())
	}
}

func assertSlicesAreEqual(t *testing.T, expected []int, actual []int, msgAndArgs ...interface{}) {
	assert.Subset(t, expected, actual, msgAndArgs...)
	assert.Subset(t, actual, expected, msgAndArgs...)
}

func assertMapsAreEqual(t *testing.T, expected IntMap, actual IntMap, msgAndArgs ...interface{}) {
	assert.Subset(t, expected.Keys(), actual.Keys())
	assert.Subset(t, actual.Keys(), expected.Keys())
	expected.Each(func(key int, val int) {
		assert.Equal(t, expected.data[key], actual.data[key], msgAndArgs...)
	})
}
