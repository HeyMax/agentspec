package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// PromptString asks for a string input with a default value.
func PromptString(label, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", Bold(label), def)
	} else {
		fmt.Printf("%s: ", Bold(label))
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

// PromptBool asks a yes/no question.
func PromptBool(label string, def bool) bool {
	defStr := "y/N"
	if def {
		defStr = "Y/n"
	}
	fmt.Printf("%s [%s]: ", Bold(label), defStr)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return def
	}
	return line == "y" || line == "yes"
}

// PromptChoice asks the user to pick from a list of options.
func PromptChoice(label string, options []string, def int) int {
	fmt.Printf("%s\n", Bold(label))
	for i, opt := range options {
		marker := "  "
		if i == def {
			marker = Green("→ ")
		}
		fmt.Printf("%s%d) %s\n", marker, i+1, opt)
	}
	fmt.Printf("Choice [%d]: ", def+1)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	var choice int
	if _, err := fmt.Sscanf(line, "%d", &choice); err != nil || choice < 1 || choice > len(options) {
		return def
	}
	return choice - 1
}

// PromptList asks for multiple comma-separated values.
func PromptList(label, example string) []string {
	fmt.Printf("%s (comma-separated, e.g. %s): ", Bold(label), Faint(example))
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	parts := strings.Split(line, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
