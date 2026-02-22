package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestCmd creates a fresh command with all flags configured.
// This ensures each test has isolated flag state.
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	setupFlags(cmd)
	return cmd
}

func TestParseFlags(t *testing.T) {
	// Save and restore environment variables
	origGOFILE := os.Getenv("GOFILE")
	origGOPACKAGE := os.Getenv("GOPACKAGE")
	t.Cleanup(func() {
		os.Setenv("GOFILE", origGOFILE)
		os.Setenv("GOPACKAGE", origGOPACKAGE)
	})

	t.Run("required type flag missing", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{})
		require.NoError(t, err)

		_, _, err = parseFlags(cmd.Flags())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty input type")
	})

	t.Run("required src falls back to GOFILE", func(t *testing.T) {
		os.Setenv("GOFILE", "myfile.go")
		os.Setenv("GOPACKAGE", "mypkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion"})
		require.NoError(t, err)

		inCfg, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		absPath, _ := filepath.Abs("myfile.go")
		assert.Equal(t, absPath, inCfg.Source)
		assert.Equal(t, "myUnion", inCfg.Type)
		assert.Equal(t, "mypkg", outCfg.OutPkg)
	})

	t.Run("missing src and GOFILE errors", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "mypkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion"})
		require.NoError(t, err)

		_, _, err = parseFlags(cmd.Flags())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "one of src or GOFILE must be set")
	})

	t.Run("missing out-pkg and GOPACKAGE errors", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion"})
		require.NoError(t, err)

		_, _, err = parseFlags(cmd.Flags())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "one of out-pkg or GOPACKAGE must be set")
	})

	t.Run("out-type defaults to capitalized type + Union", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)
		assert.Equal(t, "MyUnionUnion", outCfg.OutType)
	})

	t.Run("out-type defaults with lowercase single char", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "x"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)
		assert.Equal(t, "XUnion", outCfg.OutType)
	})

	t.Run("out-file defaults to src_gunion.go", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion", "--src", "types.go"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)
		assert.Equal(t, "types_gunion.go", outCfg.OutFile)
	})

	t.Run("out-file defaults with path", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion", "--src", "internal/types/types.go"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)
		assert.Equal(t, "internal/types/types_gunion.go", outCfg.OutFile)
	})

	t.Run("all explicit flags", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{
			"--type", "inputType",
			"--out-type", "OutputType",
			"--src", "source.go",
			"--out-file", "output.go",
			"--out-pkg", "outpkg",
			"--public-value",
			"--no-getters",
			"--no-setters",
			"--no-switch",
			"--no-default",
		})
		require.NoError(t, err)

		inCfg, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		absPath, _ := filepath.Abs("source.go")
		assert.Equal(t, config.InputConfig{
			Source: absPath,
			Type:   "inputType",
		}, inCfg)

		assert.Equal(t, config.OutputConfig{
			OutType:     "OutputType",
			OutFile:     "output.go",
			OutPkg:      "outpkg",
			PublicValue: true,
			Getters:     false,
			Setters:     false,
			Switch:      false,
			Default:     false,
		}, outCfg)
	})

	t.Run("boolean flags default to generating features", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		// By default, all features are enabled (no-* flags are false)
		assert.False(t, outCfg.PublicValue)
		assert.True(t, outCfg.Getters)
		assert.True(t, outCfg.Setters)
		assert.True(t, outCfg.Switch)
		assert.True(t, outCfg.Default)
	})

	t.Run("short flags work", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"-t", "myType", "--src", "file.go", "-o", "out.go"})
		require.NoError(t, err)

		inCfg, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		assert.Equal(t, "myType", inCfg.Type)
		assert.Equal(t, "out.go", outCfg.OutFile)
	})

	t.Run("explicit out-pkg overrides GOPACKAGE", func(t *testing.T) {
		os.Setenv("GOFILE", "test.go")
		os.Setenv("GOPACKAGE", "envpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion", "--out-pkg", "flagpkg"})
		require.NoError(t, err)

		_, outCfg, err := parseFlags(cmd.Flags())
		require.NoError(t, err)
		assert.Equal(t, "flagpkg", outCfg.OutPkg)
	})

	t.Run("explicit src overrides GOFILE", func(t *testing.T) {
		os.Setenv("GOFILE", "env.go")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion", "--src", "flag.go"})
		require.NoError(t, err)

		inCfg, _, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		absPath, _ := filepath.Abs("flag.go")
		assert.Equal(t, absPath, inCfg.Source)
	})

	t.Run("source path is converted to absolute", func(t *testing.T) {
		os.Setenv("GOFILE", "")
		os.Setenv("GOPACKAGE", "testpkg")

		cmd := newTestCmd()
		err := cmd.Flags().Parse([]string{"--type", "myUnion", "--src", "./relative/path/file.go"})
		require.NoError(t, err)

		inCfg, _, err := parseFlags(cmd.Flags())
		require.NoError(t, err)

		assert.True(t, filepath.IsAbs(inCfg.Source))
		assert.Contains(t, inCfg.Source, "relative/path/file.go")
	})
}
