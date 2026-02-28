package codegen_test

import (
	"os"
	"testing"

	"github.com/sidkurella/gunion/internal/codegen"
	"github.com/sidkurella/gunion/internal/config"
	testdata_aliasedimport "github.com/sidkurella/gunion/internal/testdata/aliasedimport"
	testdata_basic "github.com/sidkurella/gunion/internal/testdata/basic"
	testdata_externalimport "github.com/sidkurella/gunion/internal/testdata/externalimport"
	testdata_generics "github.com/sidkurella/gunion/internal/testdata/generics"
	testdata_imported "github.com/sidkurella/gunion/internal/testdata/imported"
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
				Match:       true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_basic.Representation,
			outFile:  "../testdata/basic/gen.go",
		},
		{
			name: "generics, all defaults",
			inConfig: config.OutputConfig{
				OutType:     "MyUnionUnion",
				OutPkg:      "generics",
				OutFile:     tmpDir + "/generics_gunion.go",
				PublicValue: false,
				Getters:     true,
				Setters:     true,
				Match:       true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_generics.Representation,
			outFile:  "../testdata/generics/gen.go",
		},
		{
			name: "imported, all defaults",
			inConfig: config.OutputConfig{
				OutType:     "MyUnionUnion",
				OutPkg:      "imported",
				OutFile:     tmpDir + "/imported_gunion.go",
				PublicValue: false,
				Getters:     true,
				Setters:     true,
				Match:       true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_imported.Representation,
			outFile:  "../testdata/imported/gen.go",
		},
		{
			name: "externalimport, all defaults",
			inConfig: config.OutputConfig{
				OutType:     "MyUnionUnion",
				OutPkg:      "externalimport",
				OutFile:     tmpDir + "/externalimport_gunion.go",
				PublicValue: false,
				Getters:     true,
				Setters:     true,
				Match:       true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_externalimport.Representation,
			outFile:  "../testdata/externalimport/gen.go",
		},
		{
			name: "aliasedimport, all defaults",
			inConfig: config.OutputConfig{
				OutType:     "MyUnionUnion",
				OutPkg:      "aliasedimport",
				OutFile:     tmpDir + "/aliasedimport_gunion.go",
				PublicValue: false,
				Getters:     true,
				Setters:     true,
				Match:       true,
				Default:     false,
			},
			outError: nil,
			inNamed:  testdata_aliasedimport.Representation,
			outFile:  "../testdata/aliasedimport/gen.go",
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
