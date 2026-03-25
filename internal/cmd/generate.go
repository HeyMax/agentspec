package cmd

import (
	"fmt"

	"github.com/chenzhaohao/agentspec/internal/generate"
	"github.com/chenzhaohao/agentspec/internal/util"
	"github.com/spf13/cobra"
)

var (
	genTarget    string
	genOutputDir string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code scaffolding from spec",
	Long:  `Generate backend, frontend, or test code from an AgentSpec file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if genTarget == "" {
			return fmt.Errorf("--target is required (valid: backend, frontend, test)")
		}

		specPath := GetSpecFile()
		outputDir := genOutputDir
		if outputDir == "" {
			outputDir = "./generated/" + genTarget
		}

		util.PrintHeader(fmt.Sprintf("Generating %s code from %s", genTarget, specPath))

		err := generate.Run(generate.GenerateOptions{
			SpecPath:  specPath,
			Target:    generate.Target(genTarget),
			OutputDir: outputDir,
		})
		if err != nil {
			util.PrintFail(err.Error())
			return err
		}

		util.PrintOK(fmt.Sprintf("Generated %s code in %s", genTarget, outputDir))
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&genTarget, "target", "t", "", "generation target: backend, frontend, test")
	generateCmd.Flags().StringVarP(&genOutputDir, "output", "o", "", "output directory (default: ./generated/<target>)")
	rootCmd.AddCommand(generateCmd)
}
