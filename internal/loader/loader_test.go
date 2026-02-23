package loader_test

import (
	"testing"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/sidkurella/gunion/internal/loader"
	"github.com/sidkurella/gunion/internal/testdata/aliasedimport"
	"github.com/sidkurella/gunion/internal/testdata/basic"
	"github.com/sidkurella/gunion/internal/testdata/externalimport"
	"github.com/sidkurella/gunion/internal/testdata/generics"
	"github.com/sidkurella/gunion/internal/testdata/imported"
	"github.com/sidkurella/gunion/internal/testdata/torture"
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
				Source: "../testdata/basic/basic.go",
				Type:   "myUnion",
			},
			outNamed: basic.Representation,
		},
		{
			name: "imported",
			inConfig: config.InputConfig{
				Source: "../testdata/imported/imported.go",
				Type:   "myUnion",
			},
			outNamed: imported.Representation,
		},
		{
			name: "externalimport",
			inConfig: config.InputConfig{
				Source: "../testdata/externalimport/externalimport.go",
				Type:   "myUnion",
			},
			outNamed: externalimport.Representation,
		},
		{
			name: "generics",
			inConfig: config.InputConfig{
				Source: "../testdata/generics/generics.go",
				Type:   "myUnion",
			},
			outNamed: generics.Representation,
		},
		{
			name: "torture",
			inConfig: config.InputConfig{
				Source: "../testdata/torture/torture.go",
				Type:   "myUnion",
			},
			outNamed: torture.Representation,
		},
		{
			name: "aliasedimport",
			inConfig: config.InputConfig{
				Source: "../testdata/aliasedimport/aliasedimport.go",
				Type:   "myUnion",
			},
			outNamed: aliasedimport.Representation,
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
