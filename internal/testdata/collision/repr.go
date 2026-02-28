package collision

import "github.com/sidkurella/gunion/internal/types"

// Representation is the parsed type representation of myUnion.
// This struct has fields named _variant, _inner, and Invalid, which collide
// with the default generated names. The generator should use __variant,
// __inner, and _Invalid instead.
var Representation = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/testdata/collision",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "_variant", Type: types.Basic{Name: "int"}}},
			{Var: types.Var{Name: "_inner", Type: types.Basic{Name: "string"}}},
			{Var: types.Var{Name: "Invalid", Type: types.Basic{Name: "bool"}}},
			{Var: types.Var{Name: "c", Type: types.Basic{Name: "float64"}}},
		},
	},
}
