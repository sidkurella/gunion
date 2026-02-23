package torture

import "github.com/sidkurella/gunion/internal/types"

// Representation is the parsed type representation of myUnion.
// This torture test covers many edge cases for type parsing.
var Representation = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/testdata/torture",
	Type: types.Struct{
		Fields: []types.Field{
			// Basic types
			{Var: types.Var{Name: "a", Type: types.Basic{Name: "int"}}},
			{Var: types.Var{Name: "b", Type: types.Basic{Name: "string"}}},

			// Pointer
			{Var: types.Var{Name: "c", Type: types.Pointer{Elem: types.Basic{Name: "float64"}}}},

			// Slice of pointers
			{Var: types.Var{Name: "d", Type: types.Slice{Elem: types.Pointer{Elem: types.Basic{Name: "int"}}}}},

			// Array
			{Var: types.Var{Name: "e", Type: types.Array{Len: 5, Elem: types.Basic{Name: "byte"}}}},

			// Nested maps
			{Var: types.Var{Name: "f", Type: types.Map{Key: types.Basic{Name: "string"}, Value: types.Slice{Elem: types.Basic{Name: "int"}}}}},
			{Var: types.Var{Name: "g", Type: types.Map{Key: types.Basic{Name: "int"}, Value: types.Map{Key: types.Basic{Name: "string"}, Value: types.Basic{Name: "bool"}}}}},

			// Channels (SendRecv is zero value, so elided)
			{Var: types.Var{Name: "h", Type: types.Chan{Elem: types.Basic{Name: "int"}}}},
			{Var: types.Var{Name: "i", Type: types.Chan{Direction: types.RecvOnly, Elem: types.Basic{Name: "string"}}}},
			{Var: types.Var{Name: "j", Type: types.Chan{Direction: types.SendOnly, Elem: types.Basic{Name: "bool"}}}},

			// Empty function
			{Var: types.Var{Name: "k", Type: types.Signature{}}},

			// Function with params and returns
			{Var: types.Var{Name: "l", Type: types.Signature{
				Params:  []types.Var{{Name: "a", Type: types.Basic{Name: "int"}}, {Name: "b", Type: types.Basic{Name: "string"}}},
				Returns: []types.Var{{Type: types.Basic{Name: "int"}}, {Type: types.Named{Name: "error"}}},
			}}},

			// Variadic function
			{Var: types.Var{Name: "m", Type: types.Signature{
				Params:   []types.Var{{Name: "format", Type: types.Basic{Name: "string"}}, {Name: "args", Type: types.Slice{Elem: types.Named{Name: "any"}}}},
				Returns:  []types.Var{{Type: types.Basic{Name: "string"}}},
				Variadic: true,
			}}},

			// Named params and returns
			{Var: types.Var{Name: "n", Type: types.Signature{
				Params:  []types.Var{{Name: "x", Type: types.Basic{Name: "int"}}, {Name: "y", Type: types.Basic{Name: "int"}}},
				Returns: []types.Var{{Name: "sum", Type: types.Basic{Name: "int"}}, {Name: "diff", Type: types.Basic{Name: "int"}}},
			}}},

			// Function returning function
			{Var: types.Var{Name: "o", Type: types.Signature{
				Params: []types.Var{{Name: "multiplier", Type: types.Basic{Name: "int"}}},
				Returns: []types.Var{{Type: types.Signature{
					Params:  []types.Var{{Type: types.Basic{Name: "int"}}},
					Returns: []types.Var{{Type: types.Basic{Name: "int"}}},
				}}},
			}}},

			// Function with callback param
			{Var: types.Var{Name: "p", Type: types.Signature{
				Params: []types.Var{{Name: "callback", Type: types.Signature{
					Params:  []types.Var{{Type: types.Basic{Name: "int"}}},
					Returns: []types.Var{{Type: types.Basic{Name: "bool"}}},
				}}},
				Returns: []types.Var{{Type: types.Named{Name: "error"}}},
			}}},

			// Function with channel params
			{Var: types.Var{Name: "q", Type: types.Signature{
				Params: []types.Var{
					{Name: "in", Type: types.Chan{Direction: types.RecvOnly, Elem: types.Basic{Name: "int"}}},
					{Name: "out", Type: types.Chan{Direction: types.SendOnly, Elem: types.Basic{Name: "int"}}},
				},
			}}},

			// Function with external types
			{Var: types.Var{Name: "r", Type: types.Signature{
				Params:  []types.Var{{Name: "ctx", Type: types.Named{Name: "Context", Package: "context"}}},
				Returns: []types.Var{{Type: types.Named{Name: "error"}}},
			}}},
			{Var: types.Var{Name: "s", Type: types.Signature{
				Params:  []types.Var{{Name: "w", Type: types.Named{Name: "Writer", Package: "io"}}, {Name: "r", Type: types.Named{Name: "Reader", Package: "io"}}},
				Returns: []types.Var{{Type: types.Basic{Name: "int64"}}, {Type: types.Named{Name: "error"}}},
			}}},

			// Pointer to function
			{Var: types.Var{Name: "t", Type: types.Pointer{Elem: types.Signature{
				Params:  []types.Var{{Type: types.Basic{Name: "int"}}},
				Returns: []types.Var{{Type: types.Basic{Name: "int"}}},
			}}}},

			// Slice of functions
			{Var: types.Var{Name: "u", Type: types.Slice{Elem: types.Signature{
				Returns: []types.Var{{Type: types.Named{Name: "error"}}},
			}}}},

			// Map with function values
			{Var: types.Var{Name: "v", Type: types.Map{
				Key: types.Basic{Name: "string"},
				Value: types.Signature{
					Params:  []types.Var{{Type: types.Basic{Name: "int"}}},
					Returns: []types.Var{{Type: types.Basic{Name: "int"}}},
				},
			}}},

			// Triple pointer
			{Var: types.Var{Name: "w", Type: types.Pointer{Elem: types.Pointer{Elem: types.Pointer{Elem: types.Basic{Name: "int"}}}}}},

			// Nested pointer/slice/array
			{Var: types.Var{Name: "x", Type: types.Pointer{Elem: types.Slice{Elem: types.Pointer{Elem: types.Array{Len: 3, Elem: types.Basic{Name: "int"}}}}}}},

			// Function with pointer params and return
			{Var: types.Var{Name: "y", Type: types.Signature{
				Params:  []types.Var{{Type: types.Pointer{Elem: types.Basic{Name: "int"}}}, {Type: types.Pointer{Elem: types.Pointer{Elem: types.Basic{Name: "string"}}}}},
				Returns: []types.Var{{Type: types.Pointer{Elem: types.Basic{Name: "bool"}}}},
			}}},

			// any type
			{Var: types.Var{Name: "z", Type: types.Named{Name: "any"}}},
			{Var: types.Var{Name: "aa", Type: types.Slice{Elem: types.Named{Name: "any"}}}},
			{Var: types.Var{Name: "bb", Type: types.Map{Key: types.Named{Name: "any"}, Value: types.Named{Name: "any"}}}},

			// Multiple returns
			{Var: types.Var{Name: "cc", Type: types.Signature{
				Returns: []types.Var{{Type: types.Basic{Name: "int"}}, {Type: types.Basic{Name: "int"}}, {Type: types.Named{Name: "error"}}, {Type: types.Named{Name: "error"}}},
			}}},

			// Variadic only
			{Var: types.Var{Name: "dd", Type: types.Signature{
				Params:   []types.Var{{Type: types.Slice{Elem: types.Basic{Name: "string"}}}},
				Variadic: true,
			}}},

			// Deeply nested functions
			{Var: types.Var{Name: "ee", Type: types.Signature{
				Params: []types.Var{{Type: types.Signature{
					Params:  []types.Var{{Type: types.Signature{Params: []types.Var{{Type: types.Basic{Name: "int"}}}, Returns: []types.Var{{Type: types.Basic{Name: "int"}}}}}},
					Returns: []types.Var{{Type: types.Signature{Params: []types.Var{{Type: types.Basic{Name: "int"}}}, Returns: []types.Var{{Type: types.Basic{Name: "int"}}}}}},
				}}},
				Returns: []types.Var{{Type: types.Signature{
					Params:  []types.Var{{Type: types.Signature{Params: []types.Var{{Type: types.Basic{Name: "int"}}}, Returns: []types.Var{{Type: types.Basic{Name: "int"}}}}}},
					Returns: []types.Var{{Type: types.Signature{Params: []types.Var{{Type: types.Basic{Name: "int"}}}, Returns: []types.Var{{Type: types.Basic{Name: "int"}}}}}},
				}}},
			}}},

			// Alias preservation - verify alias names are used, not underlying types
			{Var: types.Var{Name: "ff", Type: types.Basic{Name: "rune"}}},  // alias for int32
			{Var: types.Var{Name: "gg", Type: types.Basic{Name: "int32"}}}, // NOT an alias
			{Var: types.Var{Name: "hh", Type: types.Basic{Name: "byte"}}},  // alias for uint8
			{Var: types.Var{Name: "ii", Type: types.Basic{Name: "uint8"}}}, // NOT an alias

			// Generic instantiation tests
			{Var: types.Var{Name: "jj", Type: types.Named{
				Name:       "Generic",
				Package:    "github.com/sidkurella/gunion/internal/testdata/torture",
				TypeParams: []types.TypeParam{{Name: "T", Constraint: types.Named{Name: "any"}}},
				TypeArgs:   []types.Type{types.Basic{Name: "int"}},
			}}},
			{Var: types.Var{Name: "kk", Type: types.Named{
				Name:       "Generic",
				Package:    "github.com/sidkurella/gunion/internal/testdata/torture",
				TypeParams: []types.TypeParam{{Name: "T", Constraint: types.Named{Name: "any"}}},
				TypeArgs:   []types.Type{types.Basic{Name: "string"}},
			}}},
			{Var: types.Var{Name: "ll", Type: types.Named{
				Name:       "Generic",
				Package:    "github.com/sidkurella/gunion/internal/testdata/torture",
				TypeParams: []types.TypeParam{{Name: "T", Constraint: types.Named{Name: "any"}}},
				TypeArgs:   []types.Type{types.Pointer{Elem: types.Basic{Name: "float64"}}},
			}}},
			{Var: types.Var{Name: "mm", Type: types.Named{
				Name:    "TwoParam",
				Package: "github.com/sidkurella/gunion/internal/testdata/torture",
				TypeParams: []types.TypeParam{
					{Name: "K", Constraint: types.Named{Name: "comparable"}},
					{Name: "V", Constraint: types.Named{Name: "any"}},
				},
				TypeArgs: []types.Type{types.Basic{Name: "int"}, types.Basic{Name: "string"}},
			}}},
			{Var: types.Var{Name: "nn", Type: types.Named{
				Name:       "Generic",
				Package:    "github.com/sidkurella/gunion/internal/testdata/torture",
				TypeParams: []types.TypeParam{{Name: "T", Constraint: types.Named{Name: "any"}}},
				TypeArgs: []types.Type{types.Named{
					Name:       "Generic",
					Package:    "github.com/sidkurella/gunion/internal/testdata/torture",
					TypeParams: []types.TypeParam{{Name: "T", Constraint: types.Named{Name: "any"}}},
					TypeArgs:   []types.Type{types.Basic{Name: "int"}},
				}},
			}}},
			{Var: types.Var{Name: "oo", Type: types.Named{Name: "Map", Package: "sync"}}},
			{Var: types.Var{Name: "pp", Type: types.Named{Name: "Pool", Package: "sync"}}},
		},
	},
}
