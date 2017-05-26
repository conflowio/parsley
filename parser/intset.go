package parser

type IntSet []int

// NewIntMap creates a new int set
func NewIntSet() IntSet {
	return make([]int, 0)
}

func (i IntSet) Clone() IntSet {
	i2 := make([]int, len(i))
	copy(i2, i)
	return i2
}

func (i IntSet) Append(val int) IntSet {
	for _, v := range i {
		if v == val {
			return i
		}
	}
	return append(i, val)
}

func (i IntSet) Union(i2 IntSet) IntSet {
	i3 := i.Clone()
	for _, v := range i2 {
		i3.Append(v)
	}
	return i3
}
