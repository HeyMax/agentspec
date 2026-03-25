package cmd

import (
	"fmt"

	"github.com/chenzhaohao/agentspec/internal/check"
	"github.com/chenzhaohao/agentspec/internal/util"
	"github.com/spf13/cobra"
)

var checkSourceDir string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verify existing code conforms to spec",
	Long:  `Scan source code and verify it conforms to the spec (API coverage, model consistency).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := GetSpecFile()
		sourceDir := checkSourceDir
		if sourceDir == "" {
			sourceDir = "."
		}

		util.PrintHeader(fmt.Sprintf("Checking %s against %s", sourceDir, specPath))

		result, err := check.Check(specPath, sourceDir)
		if err != nil {
			util.PrintFail(err.Error())
			return err
		}

		for _, f := range result.Findings {
			if f.Pass {
				util.PrintOK(fmt.Sprintf("[%s] %s", f.Section, f.Detail))
			} else {
				util.PrintFail(fmt.Sprintf("[%s] %s", f.Section, f.Detail))
			}
		}

		fmt.Println()
		total := len(result.Findings)
		passed := result.PassCount()
		failed := result.FailCount()
		fmt.Printf("Results: %s passed, %s failed out of %d checks\n",
			util.Green(fmt.Sprintf("%d", passed)),
			util.Red(fmt.Sprintf("%d", failed)),
			total,
		)

		if failed > 0 {
			return fmt.Errorf("%d compliance checks failed", failed)
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVar(&checkSourceDir, "source", "", "source directory to scan (default: current directory)")
	rootCmd.AddCommand(checkCmd)
}
