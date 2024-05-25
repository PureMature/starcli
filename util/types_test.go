package util_test

import (
	"reflect"
	"testing"

	"github.com/PureMature/starcli/util"
	"go.starlark.net/starlark"
)

func TestOneOrMany_Unpack(t *testing.T) {
	type OneOrManyString = util.OneOrMany[starlark.String]
	tests := []struct {
		name     string
		target   *OneOrManyString
		inV      starlark.Value
		want     []starlark.String
		wantNull bool
		wantErr  bool
	}{
		{
			name:    "nil",
			target:  nil,
			inV:     starlark.String("Hello"),
			wantErr: true,
		},
		{
			name:    "int",
			target:  util.NewOneOrMany[starlark.String](starlark.String("")),
			inV:     starlark.MakeInt(42),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			n := tt.name
			p := tt.target
			// run
			err := starlark.UnpackArgs("test", []starlark.Value{tt.inV}, nil, "v?", p)
			// check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Nullable[%s].Unpack() error = %v, wantErr %v", n, err, tt.wantErr)
			} else if err != nil {
				t.Logf("Nullable[%s].Unpack() error = %v", n, err)
			}
			if tt.wantErr {
				return
			}
			// check methods
			if tt.wantNull != p.IsNull() {
				t.Errorf("Nullable[%s].IsNull() got = %v, want %v", n, p.IsNull(), tt.wantNull)
			}
			if !reflect.DeepEqual(p.Slice(), tt.want) {
				t.Errorf("Nullable[%s].Unpack() got = %v, want %v", n, p.Slice(), tt.want)
			}
		})
	}
}
