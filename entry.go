package btree

import (
	"cmp"
	"fmt"
)

type Entry[K cmp.Ordered, V any] interface {
	Key() K
	Value() V
}

type entry[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func (e entry[K, V]) String() string {
	return fmt.Sprintf("%v=%v", e.key, e.value)
}

func (e entry[K, V]) Key() K {
	return e.key
}

func (e entry[K, V]) Value() V {
	return e.value
}
