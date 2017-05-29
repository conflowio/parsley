package data

// IntMap is an immutable int-int map
type IntMap struct {
	data map[int]int
}

// NewIntMap creates a new IntMap instance
func NewIntMap() IntMap {
	return IntMap{make(map[int]int)}
}

func (i IntMap) clone() IntMap {
	i2 := IntMap{make(map[int]int, len(i.data))}
	for k, v := range i.data {
		i2.data[k] = v
	}
	return i2
}

// Get returns with the value for the given key
func (i IntMap) Get(key int) int {
	return i.data[key]
}

// Keys returns with the map keys
func (i IntMap) Keys() []int {
	n := 0
	keys := make([]int, len(i.data))
	for k := range i.data {
		keys[n] = k
		n++
	}
	return keys
}

// Inc increments the value in the integer map. If it doesn't exist it will be inserted
func (i IntMap) Inc(val int) IntMap {
	i2 := i.clone()
	if _, ok := i2.data[val]; !ok {
		i2.data[val] = 1
	} else {
		i2.data[val]++
	}
	return i2
}

// Filter creates a new int map including only the keys present in the given integer set
func (i IntMap) Filter(keys IntSet) IntMap {
	i2 := NewIntMap()
	keys.Each(func(key int) {
		if v, ok := i.data[key]; ok {
			i2.data[key] = v
		}
	})
	return i2
}

// Each runs the given function on all elements of the map
func (i IntMap) Each(f func(key int, val int)) {
	for k, v := range i.data {
		f(k, v)
	}
}
