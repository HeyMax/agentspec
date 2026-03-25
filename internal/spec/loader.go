package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadSpec loads and parses a FeatureSpec from a YAML file.
func LoadSpec(path string) (*FeatureSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}

	var s FeatureSpec
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	absPath, _ := filepath.Abs(path)
	s.FilePath = absPath
	s.LoadedAt = time.Now()

	if s.Kind == "" {
		s.Kind = "FeatureSpec"
	}
	if s.Version == "" {
		s.Version = "1.0"
	}

	return &s, nil
}

// LoadSpecFromBytes parses a FeatureSpec from YAML bytes.
func LoadSpecFromBytes(data []byte) (*FeatureSpec, error) {
	var s FeatureSpec
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	s.LoadedAt = time.Now()
	return &s, nil
}

// WriteSpec writes a FeatureSpec to a YAML file.
func WriteSpec(path string, s *FeatureSpec) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling spec: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
