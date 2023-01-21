/*
Copyright © 2023 Siddharth Kurella

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gunion",
	Short: "Generates tagged unions based on a struct definition",
	Long: `Generates a tagged union based on a struct definition.

The resultant union provides a variant field indicating which of the fields is valid.
The first field of the union should be your default type.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		cfg, err := parseFlags(flags)
		if err != nil {
			return err
		}

		fmt.Printf("%#v\n", cfg)
		return nil
	},
}

func parseFlags(flags *pflag.FlagSet) (config.Config, error) {
	inType, err := flags.GetString("type")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse type flag: %w", err)
	}
	if inType == "" {
		return config.Config{}, fmt.Errorf("received empty input type")
	}

	outType, err := flags.GetString("out-type")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse out-type flag: %w", err)
	}
	if outType == "" {
		outType = strings.ToUpper(inType[0:1]) + inType[1:] + "Union"
	}

	src, err := flags.GetString("src")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse src flag: %w", err)
	}
	if src == "" {
		goFile := os.Getenv("GOFILE")
		if goFile == "" {
			return config.Config{}, fmt.Errorf("one of src or GOFILE must be set")
		}
		src = goFile
	}

	outFile, err := flags.GetString("out-file")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse out-file flag: %w", err)
	}
	if outFile == "" {
		ext := filepath.Ext(src)
		basename := strings.TrimSuffix(src, ext)
		outFile = basename + "_gunion" + ext
	}

	outPkg, err := flags.GetString("out-pkg")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse out-pkg flag: %w", err)
	}
	if outPkg == "" {
		goPkg := os.Getenv("GOPACKAGE")
		if goPkg == "" {
			return config.Config{}, fmt.Errorf("one of out-pkg or GOPACKAGE must be set")
		}
		outPkg = goPkg
	}

	publicValue, err := flags.GetBool("public-value")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse public-value flag: %w", err)
	}

	getters, err := flags.GetBool("getters")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse getters flag: %w", err)
	}

	setters, err := flags.GetBool("setters")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse setters flag: %w", err)
	}

	makeSwitch, err := flags.GetBool("switch")
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse switch flag: %w", err)
	}

	return config.Config{
		Type:        inType,
		OutFile:     outFile,
		OutPkg:      outPkg,
		PublicValue: publicValue,
		Getters:     getters,
		Setters:     setters,
		Switch:      makeSwitch,
	}, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gunion.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("type", "t", "", "Type to generate the union from.")
	rootCmd.MarkFlagRequired("type")
	rootCmd.Flags().String(
		"out-type", "",
		"Output type name. If not specified, capitalizes the input type name and suffixes with Union.",
	)
	rootCmd.Flags().String(
		"src", "",
		"File to read from. If not present, populates with value from GOFILE environment variable.",
	)
	rootCmd.Flags().StringP("out-file", "o", "", "Output file name. If not specified, uses src_gunion.go")
	rootCmd.Flags().String("out-pkg", "", "Output package name. If not specified, uses current package.")
	rootCmd.Flags().Bool("public-value", false, "Directly export union value fields.")
	rootCmd.Flags().BoolP("getters", "g", true, "Generate getters for union members.")
	rootCmd.Flags().BoolP("setters", "s", true, "Generate setters for union members.")
	rootCmd.Flags().Bool("switch", true, "Generate switch function for union members.")
}
