package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/HeyMax/agentspec/internal/diff"
	"github.com/HeyMax/agentspec/internal/spec"
	"github.com/HeyMax/agentspec/internal/util"
	"github.com/spf13/cobra"
)

var diffOldFile string
var diffNewFile string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Analyze impact of spec changes",
	Long:  `Compare two versions of a spec file and show structural changes with downstream impact analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if diffOldFile == "" || diffNewFile == "" {
			return fmt.Errorf("both --old and --new flags are required")
		}

		oldSpec, err := loadSpecSource(diffOldFile)
		if err != nil {
			return fmt.Errorf("loading old spec: %w", err)
		}

		newSpec, err := loadSpecSource(diffNewFile)
		if err != nil {
			return fmt.Errorf("loading new spec: %w", err)
		}

		result := diff.Diff(oldSpec, newSpec)

		util.PrintHeader("Spec Diff")

		if !result.HasChanges() {
			util.PrintOK("No changes detected")
			return nil
		}

		fmt.Printf("%s\n\n", util.Bold(fmt.Sprintf("%d change(s) detected", len(result.Changes))))

		for _, ch := range result.Changes {
			var prefix string
			switch ch.Type {
			case diff.Added:
				prefix = util.Green(fmt.Sprintf("[+] %s", ch.Detail))
			case diff.Removed:
				prefix = util.Red(fmt.Sprintf("[-] %s", ch.Detail))
			case diff.Modified:
				prefix = util.Yellow(fmt.Sprintf("[~] %s", ch.Detail))
			}
			fmt.Println(prefix)
		}

		// Impact analysis
		impacts := diff.AnalyzeImpact(result.Changes, newSpec)

		if len(impacts) > 0 {
			fmt.Println()
			util.PrintHeader("Impact Analysis")
			fmt.Println()

			for _, imp := range impacts {
				var typeColor func(a ...interface{}) string
				switch imp.Change.Type {
				case diff.Added:
					typeColor = util.Green
				case diff.Removed:
					typeColor = util.Red
				default:
					typeColor = util.Yellow
				}

				fmt.Printf("%s %s\n", typeColor(strings.ToUpper(string(imp.Change.Type))), util.Bold(imp.Change.Path))

				if len(imp.AffectedAPIs) > 0 {
					fmt.Printf("  %s %s\n", util.Cyan("APIs:"), strings.Join(imp.AffectedAPIs, ", "))
				}
				if len(imp.AffectedPages) > 0 {
					fmt.Printf("  %s %s\n", util.Cyan("Pages:"), strings.Join(imp.AffectedPages, ", "))
				}
				if len(imp.AffectedTests) > 0 {
					fmt.Printf("  %s %s\n", util.Cyan("Tests:"), strings.Join(imp.AffectedTests, ", "))
				}
				for _, s := range imp.Suggestions {
					fmt.Printf("  %s %s\n", util.Faint("→"), s)
				}
				fmt.Println()
			}
		}

		fmt.Printf("%s %d change(s), %d downstream impact(s)\n",
			util.Bold("Summary:"),
			len(result.Changes),
			len(impacts),
		)

		return nil
	},
}

func loadSpecSource(source string) (*spec.FeatureSpec, error) {
	if strings.HasPrefix(source, "git:") {
		ref := strings.TrimPrefix(source, "git:")
		out, err := exec.Command("git", "show", ref).Output()
		if err != nil {
			return nil, fmt.Errorf("git show %s: %w", ref, err)
		}
		return spec.LoadSpecFromBytes(out)
	}
	return spec.LoadSpec(source)
}

func init() {
	diffCmd.Flags().StringVar(&diffOldFile, "old", "", "old spec file (or git:<ref>:<path>)")
	diffCmd.Flags().StringVar(&diffNewFile, "new", "", "new spec file (or git:<ref>:<path>)")
	rootCmd.AddCommand(diffCmd)
}
