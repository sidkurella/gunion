package generics

import "github.com/sidkurella/gunion/internal/types"

// Expected is the expected output from parsing myUnion.
var Expected = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/loader/testdata/generics",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "a", Type: types.Named{Name: "T", Package: "github.com/sidkurella/gunion/internal/loader/testdata/generics"}}},
			{Var: types.Var{Name: "b", Type: types.Named{Name: "U", Package: "github.com/sidkurella/gunion/internal/loader/testdata/generics"}}},
			{Var: types.Var{Name: "c", Type: types.Named{Name: "V", Package: "github.com/sidkurella/gunion/internal/loader/testdata/generics"}}},
		},
	},
	TypeParams: []types.TypeParam{
		{Name: "T", Constraint: types.Named{Name: "any"}},        // `any` is a type alias for interface{}
		{Name: "U", Constraint: types.Named{Name: "comparable"}}, // `comparable` is a built-in constraint
		{Name: "V", Constraint: types.Named{Name: "Writer", Package: "io"}},
	},
}
