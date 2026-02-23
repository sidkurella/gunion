package aliasedimport

import "github.com/sidkurella/gunion/internal/types"

// Representation is the parsed type representation of myUnion.
// Import aliases (ctx, ioalias) are NOT preserved - we use canonical package paths.
// This is correct because:
// 1. Package paths uniquely identify packages regardless of local aliases
// 2. Import aliases are scoped to the importing file only
// 3. Code generation needs canonical paths for import statements
var Representation = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/testdata/aliasedimport",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "a", Type: types.Named{Name: "Context", Package: "context"}}},
			{Var: types.Var{Name: "b", Type: types.Named{Name: "Writer", Package: "io"}}},
			{Var: types.Var{Name: "c", Type: types.Named{Name: "Reader", Package: "io"}}},
		},
	},
}
