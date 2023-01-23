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
