package util

import (
	"fmt"

	"go.starlark.net/starlark"
)

// OneOrMany encapsulates a slice of a specific type with optional default values.
type OneOrMany[T starlark.Value] struct {
	values       []T
	defaultValue T
}

// NewOneOrMany creates and returns a new OneOrMany with the given default value.
func NewOneOrMany[T starlark.Value](defaultValue T) *OneOrMany[T] {
	return &OneOrMany[T]{values: nil, defaultValue: defaultValue}
}

// Unpack implements the starlark.Value interface.
func (o *OneOrMany[T]) Unpack(v starlark.Value) error {
	if o == nil {
		return fmt.Errorf("nil OneOrMany pointer")
	}
	if _, ok := v.(starlark.NoneType); ok {
		o.values = nil
	} else if t, ok := v.(T); ok {
		o.values = []T{t}
	} else if l, ok := v.(starlark.Iterable); ok {
		sl := make([]T, 0)
		iter := l.Iterate()
		defer iter.Done()
		var x starlark.Value
		for iter.Next(&x) {
			if t, ok := x.(T); ok {
				sl = append(sl, t)
			} else {
				return fmt.Errorf("expected %T, got %s", o, x.Type())
			}
		}
		o.values = sl
	} else {
		return fmt.Errorf("expected %T or Iterable or None, got %s", o, v.Type())
	}
	return nil
}

// IsNull checks if the OneOrMany is considered null (having no values).
func (o *OneOrMany[T]) IsNull() bool {
	return o == nil || len(o.values) == 0
}

// Len returns the number of elements in OneOrMany.
func (o *OneOrMany[T]) Len() int {
	if o.IsNull() {
		return 0
	}
	return len(o.values)
}

// Slice returns the underlying slice of values.
func (o *OneOrMany[T]) Slice() []T {
	if o.IsNull() {
		return []T{}
	}
	return o.values
}

// First returns the first value if present, or the default value if not.
func (o *OneOrMany[T]) First() T {
	if len(o.values) == 0 {
		return o.defaultValue
	}
	return o.values[0]
}
