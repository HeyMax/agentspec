package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

// Target represents a code generation target.
type Target string

const (
	TargetBackend  Target = "backend"
	TargetFrontend Target = "frontend"
	TargetTest     Target = "test"
)

// GenerateOptions configures the code generation.
type GenerateOptions struct {
	SpecPath  string
	Target    Target
	OutputDir string
}

// Run loads the spec and dispatches to the appropriate generator.
func Run(opts GenerateOptions) error {
	s, err := spec.LoadSpec(opts.SpecPath)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	switch opts.Target {
	case TargetBackend:
		return generateBackend(s, opts.OutputDir)
	case TargetFrontend:
		return generateFrontend(s, opts.OutputDir)
	case TargetTest:
		return generateTests(s, opts.OutputDir)
	default:
		return fmt.Errorf("unknown target %q, valid targets: backend, frontend, test", opts.Target)
	}
}

// writeFile writes content to a file within the output directory.
func writeFile(outputDir, relPath, content string) error {
	fullPath := filepath.Join(outputDir, relPath)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}
	return os.WriteFile(fullPath, []byte(content), 0644)
}
