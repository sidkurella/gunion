package loader_test

import (
	"testing"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/sidkurella/gunion/internal/loader"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	type testcase struct {
		name     string
		inConfig config.InputConfig
		outUnion loader.Union
		outError error
	}
	cases := []testcase{
		{
			name: "basic",
			inConfig: config.InputConfig{
				Source: "testdata/basic/basic.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{
				Type: loader.Type{
					Name:   "myUnion",
					Source: "github.com/sidkurella/gunion/internal/loader/testdata/basic",
				},
				Variants: []loader.Variant{
					{
						Name: "a",
						Type: loader.Type{
							Name: "int",
						},
					},
					{
						Name: "b",
						Type: loader.Type{
							Name: "string",
						},
					},
				},
			},
			outError: nil,
		},
		{
			name: "imported",
			inConfig: config.InputConfig{
				Source: "testdata/imported/imported.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{
				Type: loader.Type{
					Name:   "myUnion",
					Source: "github.com/sidkurella/gunion/internal/loader/testdata/imported",
				},
				Variants: []loader.Variant{
					{
						Name: "a",
						Type: loader.Type{
							Name: "int",
						},
					},
					{
						Name: "b",
						Type: loader.Type{
							Name:          "MyValue",
							IndirectCount: 0,
							Source:        "github.com/sidkurella/gunion/internal/loader/testdata/imported/inner",
						},
					},
				},
			},
			outError: nil,
		},
		{
			name: "externalimport",
			inConfig: config.InputConfig{
				Source: "testdata/externalimport/externalimport.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{
				Type: loader.Type{
					Name:   "myUnion",
					Source: "github.com/sidkurella/gunion/internal/loader/testdata/externalimport",
				},
				Variants: []loader.Variant{
					{
						Name: "a",
						Type: loader.Type{
							Name: "int",
						},
					},
					{
						Name: "b",
						Type: loader.Type{
							Name:          "Package",
							IndirectCount: 1,
							Source:        "golang.org/x/tools/go/packages",
						},
					},
					{
						Name: "c",
						Type: loader.Type{
							Name:          "Context",
							IndirectCount: 0,
							Source:        "context",
						},
					},
				},
			},
			outError: nil,
		},
		{
			name: "generics",
			inConfig: config.InputConfig{
				Source: "testdata/generics/generics.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{
				Type: loader.Type{
					Name:   "myUnion",
					Source: "github.com/sidkurella/gunion/internal/loader/testdata/generics",
					TypeParams: []loader.TypeParam{
						{
							Name: "T",
							Type: loader.Type{
								Name: "any",
							},
						},
						{
							Name: "U",
							Type: loader.Type{
								Name: "comparable",
							},
						},
						{
							Name: "V",
							Type: loader.Type{
								Name:   "Writer",
								Source: "io",
							},
						},
					},
				},
				Variants: []loader.Variant{
					{
						Name: "a",
						Type: loader.Type{
							Name: "T",
						},
					},
					{
						Name: "b",
						Type: loader.Type{
							Name: "U",
						},
					},
					{
						Name: "c",
						Type: loader.Type{
							Name: "V",
						},
					},
				},
			},
			outError: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := loader.NewLoader(tc.inConfig)
			u, err := l.Load()
			if tc.outError != nil {
				require.EqualError(t, err, tc.outError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.outUnion, u)
			}
		})
	}
}
