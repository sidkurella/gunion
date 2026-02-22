package imported

import "github.com/sidkurella/gunion/internal/types"

// Expected is the expected output from parsing myUnion.
var Expected = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/loader/testdata/imported",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "a", Type: types.Basic{Name: "int"}}},
			{Var: types.Var{
				Name: "b",
				Type: types.Named{
					Name:    "MyValue",
					Package: "github.com/sidkurella/gunion/internal/loader/testdata/imported/inner",
					// Type is nil because we don't expand nested Named types
				},
			}},
		},
	},
}
