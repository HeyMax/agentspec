package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var specFile string

var rootCmd = &cobra.Command{
	Use:   "agentspec",
	Short: "The contract layer for AI-powered software development",
	Long: `AgentSpec — Let multiple AI agents collaborate on a single, structured specification.

Define your feature requirements as a machine-consumable YAML spec.
Each role (backend, frontend, QA) gets a tailored context from the same source of truth.

Commands:
  init       Create a new spec file interactively
  validate   Check spec for errors and inconsistencies
  generate   Generate code scaffolding from spec
  check      Verify existing code conforms to spec
  diff       Analyze impact of spec changes
  context    Generate role-specific AI agent context`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&specFile, "file", "f", "", "path to spec YAML file")
}

// GetSpecFile returns the spec file path, with a sensible default.
func GetSpecFile() string {
	if specFile != "" {
		return specFile
	}
	// Try common defaults
	candidates := []string{"spec.yaml", "spec.yml", "agentspec.yaml", "agentspec.yml"}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "spec.yaml"
}
