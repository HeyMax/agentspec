package cmd

import (
	"fmt"
	"os"
	"strings"

	agentcontext "github.com/chenzhaohao/agentspec/internal/context"
	"github.com/spf13/cobra"
)

var contextRole string
var contextOutput string

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate role-specific AI agent context",
	Long:  `Generate a Markdown context document filtered for a specific role (backend, frontend, qa, pm).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if contextRole == "" {
			return fmt.Errorf("--role is required (valid: %s)", strings.Join(agentcontext.RoleNames(), ", "))
		}

		specPath := GetSpecFile()
		md, err := agentcontext.Generate(specPath, contextRole)
		if err != nil {
			return err
		}

		if contextOutput != "" {
			if err := os.WriteFile(contextOutput, []byte(md), 0644); err != nil {
				return fmt.Errorf("writing output file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Context written to %s\n", contextOutput)
		} else {
			fmt.Print(md)
		}

		return nil
	},
}

func init() {
	contextCmd.Flags().StringVar(&contextRole, "role", "", "agent role: backend, frontend, qa, pm")
	contextCmd.Flags().StringVar(&contextOutput, "output", "", "write output to file instead of stdout")
	rootCmd.AddCommand(contextCmd)
}
