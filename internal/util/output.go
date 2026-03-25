package util

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Faint   = color.New(color.Faint).SprintFunc()
)

func PrintError(path, msg string) {
	fmt.Printf("  %s %s %s\n", Red("ERROR"), Faint(path+":"), msg)
}

func PrintWarn(path, msg string) {
	fmt.Printf("  %s  %s %s\n", Yellow("WARN"), Faint(path+":"), msg)
}

func PrintOK(msg string) {
	fmt.Printf("  %s %s\n", Green("✓"), msg)
}

func PrintFail(msg string) {
	fmt.Printf("  %s %s\n", Red("✗"), msg)
}

func PrintInfo(msg string) {
	fmt.Printf("  %s %s\n", Blue("ℹ"), msg)
}

func PrintHeader(title string) {
	fmt.Printf("\n%s\n", Bold(title))
	fmt.Println(Faint("─────────────────────────────────────────"))
}

func PrintSummary(errors, warnings int) {
	if errors == 0 && warnings == 0 {
		fmt.Printf("\n%s\n", Green("✓ Spec is valid — no issues found"))
	} else {
		parts := ""
		if errors > 0 {
			parts += fmt.Sprintf("%s", Red(fmt.Sprintf("%d error(s)", errors)))
		}
		if warnings > 0 {
			if parts != "" {
				parts += ", "
			}
			parts += fmt.Sprintf("%s", Yellow(fmt.Sprintf("%d warning(s)", warnings)))
		}
		fmt.Printf("\n%s\n", parts)
	}
}
