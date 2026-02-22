package loader_test

import (
	"testing"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/sidkurella/gunion/internal/loader"
	"github.com/sidkurella/gunion/internal/loader/testdata/basic"
	"github.com/sidkurella/gunion/internal/loader/testdata/externalimport"
	"github.com/sidkurella/gunion/internal/loader/testdata/generics"
	"github.com/sidkurella/gunion/internal/loader/testdata/imported"
	"github.com/sidkurella/gunion/internal/loader/testdata/torture"
	"github.com/sidkurella/gunion/internal/types"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	type testcase struct {
		name     string
		inConfig config.InputConfig
		outNamed types.Named
		outError error
	}
	cases := []testcase{
		{
			name: "basic",
			inConfig: config.InputConfig{
				Source: "testdata/basic/basic.go",
				Type:   "myUnion",
			},
			outNamed: basic.Expected,
		},
		{
			name: "imported",
			inConfig: config.InputConfig{
				Source: "testdata/imported/imported.go",
				Type:   "myUnion",
			},
			outNamed: imported.Expected,
		},
		{
			name: "externalimport",
			inConfig: config.InputConfig{
				Source: "testdata/externalimport/externalimport.go",
				Type:   "myUnion",
			},
			outNamed: externalimport.Expected,
		},
		{
			name: "generics",
			inConfig: config.InputConfig{
				Source: "testdata/generics/generics.go",
				Type:   "myUnion",
			},
			outNamed: generics.Expected,
		},
		{
			name: "torture",
			inConfig: config.InputConfig{
				Source: "testdata/torture/torture.go",
				Type:   "myUnion",
			},
			outNamed: torture.Expected,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := loader.NewLoader(tc.inConfig)
			named, err := l.Load()
			if tc.outError != nil {
				require.EqualError(t, err, tc.outError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.outNamed, named)
			}
		})
	}
}
