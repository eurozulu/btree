package btree

// helper functions for slices

func IndexOf[V comparable](v V, s []V) int {
	for i, e := range s {
		if v == e {
			return i
		}
	}
	return -1
}

func InsertAtIndex[V any](v V, s []V, index int) []V {
	var ns []V
	if index > 0 {
		ns = append(ns, s[:index]...)
	}
	ns = append(ns, v)
	if index < len(s) {
		ns = append(ns, s[index:]...)
	}
	return ns
}

func RemoveAtIndex[V any](s []V, index int) []V {
	var ns []V
	if index > 0 {
		ns = append(ns, s[:index]...)
	}
	if index+1 < len(s) {
		ns = append(ns, s[index+1:]...)
	}
	return ns
}

type stackSlice[V any] []V

func (s *stackSlice[V]) Push(v V) {
	*s = append(*s, v)
}

func (s *stackSlice[V]) Pop() (V, bool) {
	var v V
	if s.IsEmpty() {
		return v, false
	}
	last := len(*s) - 1
	v = (*s)[last]
	*s = (*s)[:last]
	return v, true
}

func (s *stackSlice[V]) Peek() (V, bool) {
	var v V
	if s.IsEmpty() {
		return v, false
	}
	last := len(*s) - 1
	return (*s)[last], true
}

func (s *stackSlice[V]) IsEmpty() bool {
	return len(*s) == 0
}
