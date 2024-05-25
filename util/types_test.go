package util

import (
	"reflect"
	"testing"

	"go.starlark.net/starlark"
)

var (
	none = starlark.None
)

func TestOneOrMany_Unpack(t *testing.T) {
	tests := []struct {
		name     string
		target   *OneOrMany[starlark.Int]
		inV      starlark.Value
		want     []starlark.Int
		wantNull bool
		wantErr  bool
	}{
		{
			name:    "nil receiver",
			target:  nil,
			inV:     starlark.MakeInt(10),
			wantErr: true,
		},
		{
			name:     "int value",
			target:   NewOneOrManyNoDefault[starlark.Int](),
			inV:      starlark.MakeInt(10),
			want:     []starlark.Int{starlark.MakeInt(10)},
			wantNull: false,
		},
		{
			name:     "none value",
			target:   NewOneOrManyNoDefault[starlark.Int](),
			inV:      none,
			want:     []starlark.Int{},
			wantNull: true,
		},
		{
			name:     "iterable value",
			target:   NewOneOrManyNoDefault[starlark.Int](),
			inV:      starlark.NewList([]starlark.Value{starlark.MakeInt(10), starlark.MakeInt(20)}),
			want:     []starlark.Int{starlark.MakeInt(10), starlark.MakeInt(20)},
			wantNull: false,
		},
		{
			name:    "wrong type value",
			target:  NewOneOrManyNoDefault[starlark.Int](),
			inV:     starlark.String("foo"),
			wantErr: true,
		},
		{
			name:    "iterable with wrong type",
			target:  NewOneOrManyNoDefault[starlark.Int](),
			inV:     starlark.NewList([]starlark.Value{starlark.MakeInt(10), starlark.String("foo")}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.target.Unpack(tt.inV)
			if (err != nil) != tt.wantErr {
				t.Errorf("OneOrMany[%s].Unpack() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			} else if err != nil {
				t.Logf("OneOrMany[%s].Unpack() error = %v", tt.name, err)
			}
			if tt.wantErr {
				return
			}
			if tt.wantNull != tt.target.IsNull() {
				t.Errorf("OneOrMany[%s].IsNull() got = %v, want %v", tt.name, tt.target.IsNull(), tt.wantNull)
			}
			if !reflect.DeepEqual(tt.target.Slice(), tt.want) {
				t.Errorf("OneOrMany[%s].Unpack() got = %v, want %v", tt.name, tt.target.Slice(), tt.want)
			}
		})
	}
}

func TestOneOrMany_First(t *testing.T) {
	tests := []struct {
		name     string
		target   *OneOrMany[starlark.Int]
		want     starlark.Int
		wantNull bool
	}{
		{
			name:     "empty no default",
			target:   NewOneOrManyNoDefault[starlark.Int](),
			want:     starlark.Int{},
			wantNull: true,
		},
		{
			name:     "empty with default",
			target:   NewOneOrMany(starlark.MakeInt(5)),
			want:     starlark.MakeInt(5),
			wantNull: true,
		},
		{
			name:     "single value",
			target:   &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:     starlark.MakeInt(10),
			wantNull: false,
		},
		{
			name:     "multiple values",
			target:   &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10), starlark.MakeInt(20)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:     starlark.MakeInt(10),
			wantNull: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OneOrMany[%s].First() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestOneOrMany_Slice(t *testing.T) {
	tests := []struct {
		name   string
		target *OneOrMany[starlark.Int]
		want   []starlark.Int
	}{
		{
			name:   "empty no default",
			target: NewOneOrManyNoDefault[starlark.Int](),
			want:   []starlark.Int{},
		},
		{
			name:   "empty with default",
			target: NewOneOrMany(starlark.MakeInt(5)),
			want:   []starlark.Int{starlark.MakeInt(5)},
		},
		{
			name:   "single value",
			target: &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:   []starlark.Int{starlark.MakeInt(10)},
		},
		{
			name:   "multiple values",
			target: &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10), starlark.MakeInt(20)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:   []starlark.Int{starlark.MakeInt(10), starlark.MakeInt(20)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.Slice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OneOrMany[%s].Slice() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestOneOrMany_IsNull(t *testing.T) {
	tests := []struct {
		name   string
		target *OneOrMany[starlark.Int]
		want   bool
	}{
		{
			name:   "nil receiver",
			target: nil,
			want:   true,
		},
		{
			name:   "empty no default",
			target: NewOneOrManyNoDefault[starlark.Int](),
			want:   true,
		},
		{
			name:   "empty with default",
			target: NewOneOrMany(starlark.MakeInt(5)),
			want:   true,
		},
		{
			name:   "single value",
			target: &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:   false,
		},
		{
			name:   "multiple values",
			target: &OneOrMany[starlark.Int]{values: []starlark.Int{starlark.MakeInt(10), starlark.MakeInt(20)}, hasDefault: true, defaultValue: starlark.MakeInt(5)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.IsNull(); got != tt.want {
				t.Errorf("OneOrMany[%s].IsNull() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
