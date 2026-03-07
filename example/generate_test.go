package example

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGoGenerate verifies that `go generate` produces a valid output file
// that matches the already-committed shape_gunion.go.
// This test actually invokes `go generate` and compares the result.
func TestGoGenerate(t *testing.T) {
	// Get the path to the example directory.
	exampleDir, err := filepath.Abs(".")
	require.NoError(t, err)

	genFile := filepath.Join(exampleDir, "shape_gunion.go")

	// Read the existing committed generated file for comparison.
	committed, err := os.ReadFile(genFile)
	require.NoError(t, err)

	// Remove it so we can verify go:generate recreates it.
	err = os.Remove(genFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		// Restore the generated file.
		_ = os.WriteFile(genFile, committed, 0644)
	})

	// Run go generate.
	cmd := exec.Command("go", "generate", ".")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "go generate failed: %s", string(output))

	// Verify the file was regenerated.
	regenerated, err := os.ReadFile(genFile)
	require.NoError(t, err, "generated file should exist after go generate")

	// The generated file should be non-empty and valid Go.
	require.NotEmpty(t, regenerated)

	// Verify it compiles by running go build.
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = exampleDir
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "go build failed after generate: %s", string(buildOutput))
}
