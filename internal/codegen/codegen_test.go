package codegen_test

import (
	"os"
	"testing"

	"github.com/sidkurella/gunion/internal/codegen"
	"github.com/sidkurella/gunion/internal/config"
	testdata_basic "github.com/sidkurella/gunion/internal/testdata/basic"
	"github.com/sidkurella/gunion/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCodeGenerator(t *testing.T) {
	type testcase struct {
		name     string
		inConfig config.OutputConfig
		inNamed  types.Named
		outFile  string // Path to expected output file for comparison
		outError error
	}

	tmpDir := t.TempDir()
	cases := []testcase{
		{
			name: "basic, all defaults",
			inConfig: config.OutputConfig{
				OutType:     "MyUnionUnion",
				OutPkg:      "basic",
				OutFile:     tmpDir + "/basic_gunion.go",
				PublicValue: false,
				Getters:     true,
				Setters:     true,
				Switch:      true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_basic.Representation,
			outFile:  "../testdata/basic/gen.go",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cg := codegen.NewCodeGenerator(tc.inConfig)
			err := cg.Generate(tc.inNamed)
			if tc.outError != nil {
				require.EqualError(t, err, tc.outError.Error())
			} else {
				require.NoError(t, err)
				actual, err := os.ReadFile(tc.inConfig.OutFile)
				require.NoError(t, err)
				expected, err := os.ReadFile(tc.outFile)
				require.NoError(t, err)
				require.Equal(t, string(expected), string(actual))
			}
		})
	}
}
