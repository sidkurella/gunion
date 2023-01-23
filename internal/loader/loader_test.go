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
			outUnion: loader.Union{},
			outError: nil,
		},
		{
			name: "imported",
			inConfig: config.InputConfig{
				Source: "testdata/imported/imported.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{},
			outError: nil,
		},
		{
			name: "externalimport",
			inConfig: config.InputConfig{
				Source: "testdata/externalimport/externalimport.go",
				Type:   "myUnion",
			},
			outUnion: loader.Union{},
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
