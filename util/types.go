package util

import (
	"fmt"

	"go.starlark.net/starlark"
)

// OneOrMany is a type alias for a slice of a specific type.
type OneOrMany[T starlark.Value] []T

// Unpack implements the starlark.Value interface.
func (s *OneOrMany[T]) Unpack(v starlark.Value) error {
	if _, ok := v.(starlark.NoneType); ok {
		*s = nil
	} else if t, ok := v.(T); ok {
		*s = []T{t}
	} else if l, ok := v.(starlark.Iterable); ok {
		sl := make([]T, 0, 1)
		iter := l.Iterate()
		defer iter.Done()
		// iterate over the iterable
		var x starlark.Value
		for iter.Next(&x) {
			if t, ok := x.(T); ok {
				sl = append(sl, t)
			} else {
				return fmt.Errorf("expected %T, got %s", s, x.Type())
			}
		}
		*s = sl
	} else {
		return fmt.Errorf("expected %T or Iterable or None, got %s", s, v.Type())
	}
	return nil
}

func (s *OneOrMany[T]) IsNull() bool {
	return s == nil
}

func (s *OneOrMany[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(*s)
}

func (s *OneOrMany[T]) Values() []T {
	if s == nil {
		return []T{}
	}
	return *s
}

func (s *OneOrMany[T]) First() T {
	if s == nil || len(*s) == 0 {
		var zero T
		return zero
	}
	return (*s)[0]
}
