package cmd

import (
	"fmt"
	"os"

	"github.com/HeyMax/agentspec/internal/schema"
	"github.com/HeyMax/agentspec/internal/spec"
	"github.com/HeyMax/agentspec/internal/util"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Check spec for errors and inconsistencies",
	Long:  `Validate a spec YAML file through structural, semantic, and cross-reference checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := GetSpecFile()
		if len(args) > 0 {
			specPath = args[0]
		}

		util.PrintHeader("Validating " + specPath)

		s, err := spec.LoadSpec(specPath)
		if err != nil {
			util.PrintFail(fmt.Sprintf("Failed to load spec: %s", err))
			os.Exit(1)
		}

		result := schema.Validate(s)

		if len(result.Diagnostics) == 0 {
			util.PrintSummary(0, 0)
			return nil
		}

		for _, d := range result.Diagnostics {
			if d.IsError() {
				util.PrintError(d.Path, d.Message)
			} else {
				util.PrintWarn(d.Path, d.Message)
			}
		}

		util.PrintSummary(result.ErrorCount(), result.WarnCount())

		if result.HasErrors() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
