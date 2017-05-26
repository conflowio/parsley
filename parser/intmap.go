package parser

// IntMap is a simple int set with a map
type IntMap map[int]int

// NewIntMap creates a new int map
func NewIntMap() IntMap {
	return make(map[int]int)
}

// Clone clones the int map
func (i IntMap) Clone() IntMap {
	i2 := NewIntMap()
	for k, v := range i {
		i2[k] = v
	}
	return i2
}

// Add adds a new value to the int map
func (i IntMap) Add(val int) IntMap {
	i2 := i.Clone()
	if _, ok := i2[val]; !ok {
		i2[val] = 1
	} else {
		i2[val]++
	}
	return i2
}

func (i IntMap) Filter(keys IntSet) IntMap {
	i2 := NewIntMap()
	for _, key := range keys {
		if v, ok := i[key]; ok {
			i2[key] = v
		}
	}
	return i2
}
