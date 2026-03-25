package context

import (
	"fmt"
	"os"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

// Generate loads a spec file and produces role-filtered Markdown context.
func Generate(specPath string, roleName string) (string, error) {
	role, ok := GetRole(roleName)
	if !ok {
		return "", fmt.Errorf("unknown role %q, valid roles: backend, frontend, qa, pm", roleName)
	}

	s, err := spec.LoadSpec(specPath)
	if err != nil {
		return "", fmt.Errorf("loading spec: %w", err)
	}

	md := RenderMarkdown(s, role)
	return md, nil
}

// GenerateToFile generates context and writes it to a file.
func GenerateToFile(specPath, roleName, outputPath string) error {
	md, err := Generate(specPath, roleName)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(md), 0644)
}
