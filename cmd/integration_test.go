package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	type testcase struct {
		name       string
		sourceFile string // relative to internal/testdata/
		typeName   string
		outPkg     string
		goldenFile string // relative to internal/testdata/
		extraFlags []string
	}

	cases := []testcase{
		{
			name:       "basic",
			sourceFile: "basic/basic.go",
			typeName:   "myUnion",
			outPkg:     "basic",
			goldenFile: "basic/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "generics",
			sourceFile: "generics/generics.go",
			typeName:   "myUnion",
			outPkg:     "generics",
			goldenFile: "generics/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "imported",
			sourceFile: "imported/imported.go",
			typeName:   "myUnion",
			outPkg:     "imported",
			goldenFile: "imported/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "externalimport",
			sourceFile: "externalimport/externalimport.go",
			typeName:   "myUnion",
			outPkg:     "externalimport",
			goldenFile: "externalimport/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "aliasedimport",
			sourceFile: "aliasedimport/aliasedimport.go",
			typeName:   "myUnion",
			outPkg:     "aliasedimport",
			goldenFile: "aliasedimport/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "basic/minimal",
			sourceFile: "basic/basic.go",
			typeName:   "myUnion",
			outPkg:     "basic",
			goldenFile: "basic/minimal/gen.go",
			extraFlags: []string{"--no-getters", "--no-setters", "--no-match", "--no-default"},
		},
		{
			name:       "basic/withdefault",
			sourceFile: "basic/basic.go",
			typeName:   "myUnion",
			outPkg:     "basic",
			goldenFile: "basic/withdefault/gen.go",
			// Default is true (the CLI default), so no --no-default flag needed.
		},
		{
			name:       "collision",
			sourceFile: "collision/collision.go",
			typeName:   "myUnion",
			outPkg:     "collision",
			goldenFile: "collision/gen.go",
			extraFlags: []string{"--no-default"},
		},
		{
			name:       "torture",
			sourceFile: "torture/torture.go",
			typeName:   "myUnion",
			outPkg:     "torture",
			goldenFile: "torture/gen.go",
			extraFlags: []string{"--no-default"},
		},
	}

	// Save and restore global state.
	origLoaderFactory := LoaderFactory
	origGeneratorFactory := GeneratorFactory
	origGOFILE := os.Getenv("GOFILE")
	origGOPACKAGE := os.Getenv("GOPACKAGE")
	origArgs := os.Args
	t.Cleanup(func() {
		LoaderFactory = origLoaderFactory
		GeneratorFactory = origGeneratorFactory
		os.Setenv("GOFILE", origGOFILE)
		os.Setenv("GOPACKAGE", origGOPACKAGE)
		os.Args = origArgs
	})

	// Clear env vars so the command relies solely on explicit flags.
	os.Setenv("GOFILE", "")
	os.Setenv("GOPACKAGE", "")

	tmpDir := t.TempDir()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Resolve source file to absolute path.
			srcAbs, err := filepath.Abs(filepath.Join("..", "internal", "testdata", tc.sourceFile))
			require.NoError(t, err)

			// Build a unique output file path.
			outFile := filepath.Join(tmpDir, strings.ReplaceAll(tc.name, "/", "_")+"_gunion.go")

			// Set os.Args to a canonical command string that matches the golden file header.
			// RunE reads os.Args for the Command field in the generated header.
			os.Args = []string{
				"gunion",
				"--type", tc.typeName,
				"--src", "source.go",
			}
			os.Args = append(os.Args, tc.extraFlags...)

			// SetArgs controls what cobra actually parses (with the real absolute paths).
			realArgs := []string{
				"--type", tc.typeName,
				"--src", srcAbs,
				"--out-type", "MyUnionUnion",
				"--out-pkg", tc.outPkg,
				"--out-file", outFile,
			}
			realArgs = append(realArgs, tc.extraFlags...)

			// Create a fresh command to avoid sticky flag state between tests.
			cmd := newRootCmd()
			cmd.SetArgs(realArgs)
			err = cmd.Execute()
			require.NoError(t, err)

			// Read actual output and golden file, compare byte-for-byte.
			actual, err := os.ReadFile(outFile)
			require.NoError(t, err)
			goldenPath := filepath.Join("..", "internal", "testdata", tc.goldenFile)
			expected, err := os.ReadFile(goldenPath)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(actual))
		})
	}
}

func TestIntegrationErrors(t *testing.T) {
	// Save and restore global state.
	origLoaderFactory := LoaderFactory
	origGeneratorFactory := GeneratorFactory
	origGOFILE := os.Getenv("GOFILE")
	origGOPACKAGE := os.Getenv("GOPACKAGE")
	origArgs := os.Args
	t.Cleanup(func() {
		LoaderFactory = origLoaderFactory
		GeneratorFactory = origGeneratorFactory
		os.Setenv("GOFILE", origGOFILE)
		os.Setenv("GOPACKAGE", origGOPACKAGE)
		os.Args = origArgs
	})

	os.Setenv("GOFILE", "")
	os.Setenv("GOPACKAGE", "")

	tmpDir := t.TempDir()

	t.Run("nonexistent source file", func(t *testing.T) {
		outFile := filepath.Join(tmpDir, "nonexistent_gunion.go")
		srcAbs, err := filepath.Abs(filepath.Join("..", "internal", "testdata", "nonexistent", "nonexistent.go"))
		require.NoError(t, err)

		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "myUnion",
			"--src", srcAbs,
			"--out-type", "MyUnionUnion",
			"--out-pkg", "nonexistent",
			"--out-file", outFile,
		})
		err = cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load type")
	})

	t.Run("nonexistent type in valid source", func(t *testing.T) {
		outFile := filepath.Join(tmpDir, "badtype_gunion.go")
		srcAbs, err := filepath.Abs(filepath.Join("..", "internal", "testdata", "basic", "basic.go"))
		require.NoError(t, err)

		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "doesNotExist",
			"--src", srcAbs,
			"--out-type", "MyUnionUnion",
			"--out-pkg", "basic",
			"--out-file", outFile,
		})
		err = cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not find type doesNotExist")
	})

	t.Run("non-struct type", func(t *testing.T) {
		outFile := filepath.Join(tmpDir, "nonstruct_gunion.go")
		srcAbs, err := filepath.Abs(filepath.Join("..", "internal", "testdata", "nonstruct", "nonstruct.go"))
		require.NoError(t, err)

		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "myUnion",
			"--src", srcAbs,
			"--out-type", "MyUnionUnion",
			"--out-pkg", "nonstruct",
			"--out-file", outFile,
		})
		err = cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected a struct type")
	})

	t.Run("source file with compile errors", func(t *testing.T) {
		outFile := filepath.Join(tmpDir, "compileerror_gunion.go")
		srcAbs, err := filepath.Abs(filepath.Join("..", "internal", "testdata", "compileerror", "compileerror.go"))
		require.NoError(t, err)

		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "myUnion",
			"--src", srcAbs,
			"--out-type", "MyUnionUnion",
			"--out-pkg", "compileerror",
			"--out-file", outFile,
		})
		err = cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "had errors")
	})

	t.Run("missing required type flag", func(t *testing.T) {
		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--src", "test.go",
			"--out-pkg", "testpkg",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "required flag")
	})

	t.Run("missing src and GOFILE", func(t *testing.T) {
		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "myUnion",
			"--out-pkg", "testpkg",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "one of src or GOFILE must be set")
	})

	t.Run("missing out-pkg and GOPACKAGE", func(t *testing.T) {
		cmd := newRootCmd()
		cmd.SetArgs([]string{
			"--type", "myUnion",
			"--src", "test.go",
		})
		err := cmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "one of out-pkg or GOPACKAGE must be set")
	})
}
