package codegen_test

import (
	"os"
	"strings"
	"testing"

	"github.com/sidkurella/gunion/internal/codegen"
	"github.com/sidkurella/gunion/internal/config"
	testdata_aliasedimport "github.com/sidkurella/gunion/internal/testdata/aliasedimport"
	testdata_basic "github.com/sidkurella/gunion/internal/testdata/basic"
	testdata_collision "github.com/sidkurella/gunion/internal/testdata/collision"
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
				OutType: "MyUnionUnion",
				OutPkg:  "basic",
				OutFile: tmpDir + "/basic_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_basic.Representation,
			outFile:  "../testdata/basic/gen.go",
		},
		{
			name: "generics, all defaults",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "generics",
				OutFile: tmpDir + "/generics_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_generics.Representation,
			outFile:  "../testdata/generics/gen.go",
		},
		{
			name: "imported, all defaults",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "imported",
				OutFile: tmpDir + "/imported_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_imported.Representation,
			outFile:  "../testdata/imported/gen.go",
		},
		{
			name: "externalimport, all defaults",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "externalimport",
				OutFile: tmpDir + "/externalimport_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_externalimport.Representation,
			outFile:  "../testdata/externalimport/gen.go",
		},
		{
			name: "aliasedimport, all defaults",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "aliasedimport",
				OutFile: tmpDir + "/aliasedimport_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_aliasedimport.Representation,
			outFile:  "../testdata/aliasedimport/gen.go",
		},
		{
			name: "basic, no getters/setters/match",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "basic",
				OutFile: tmpDir + "/basic_minimal_gunion.go",
				Getters: false,
				Setters: false,
				Match:   false,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_basic.Representation,
			outFile:  "../testdata/basic/minimal/gen.go",
		},
		{
			name: "basic, with default (no invalid variant)",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "basic",
				OutFile: tmpDir + "/basic_default_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: true,
			},
			outError: nil,
			inNamed:  testdata_basic.Representation,
			outFile:  "../testdata/basic/withdefault/gen.go",
		},
		{
			name: "collision, fields named _variant and _inner",
			inConfig: config.OutputConfig{
				OutType: "MyUnionUnion",
				OutPkg:  "collision",
				OutFile: tmpDir + "/collision_gunion.go",
				Getters: true,
				Setters: true,
				Match:   true,
				Default: false,
			},
			outError: nil,
			inNamed:  testdata_collision.Representation,
			outFile:  "../testdata/collision/gen.go",
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

func TestCodeGeneratorCommand(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("command is included in header", func(t *testing.T) {
		cfg := config.OutputConfig{
			OutType: "MyUnionUnion",
			OutPkg:  "basic",
			OutFile: tmpDir + "/basic_cmd_gunion.go",
			Command: "gunion --type myUnion --src basic.go",
			Getters: true,
			Setters: true,
			Match:   true,
		}
		cg := codegen.NewCodeGenerator(cfg)
		err := cg.Generate(testdata_basic.Representation)
		require.NoError(t, err)

		actual, err := os.ReadFile(cfg.OutFile)
		require.NoError(t, err)
		firstLine := strings.SplitN(string(actual), "\n", 2)[0]
		require.Equal(t, "// Code generated by gunion via `gunion --type myUnion --src basic.go`. DO NOT EDIT.", firstLine)
	})

	t.Run("empty command omits invocation", func(t *testing.T) {
		cfg := config.OutputConfig{
			OutType: "MyUnionUnion",
			OutPkg:  "basic",
			OutFile: tmpDir + "/basic_nocmd_gunion.go",
			Getters: true,
			Setters: true,
			Match:   true,
		}
		cg := codegen.NewCodeGenerator(cfg)
		err := cg.Generate(testdata_basic.Representation)
		require.NoError(t, err)

		actual, err := os.ReadFile(cfg.OutFile)
		require.NoError(t, err)
		firstLine := strings.SplitN(string(actual), "\n", 2)[0]
		require.Equal(t, "// Code generated by gunion. DO NOT EDIT.", firstLine)
	})
}
