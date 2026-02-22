package externalimport

import "github.com/sidkurella/gunion/internal/types"

// Expected is the expected output from parsing myUnion.
// Note: We only verify the top-level structure for external imports,
// since the full type tree of external packages is large and fragile.
var Expected = types.Named{
	Name:    "myUnion",
	Package: "github.com/sidkurella/gunion/internal/loader/testdata/externalimport",
	Type: types.Struct{
		Fields: []types.Field{
			{Var: types.Var{Name: "a", Type: types.Basic{Name: "int"}}},
			{Var: types.Var{
				Name: "b",
				Type: types.Pointer{Elem: types.Named{Name: "Package", Package: "golang.org/x/tools/go/packages"}},
			}},
			{Var: types.Var{
				Name: "c",
				Type: types.Named{Name: "Context", Package: "context"},
			}},
		},
	},
}
