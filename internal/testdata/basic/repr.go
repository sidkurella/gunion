package basic

import "github.com/sidkurella/gunion/internal/types"

// Representation is the parsed type representation of myUnion.
var Representation = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/testdata/basic",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "a", Type: types.Basic{Name: "int"}}},
			{Var: types.Var{Name: "b", Type: types.Basic{Name: "string"}}},
		},
	},
}
