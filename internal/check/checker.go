package check

import (
	"fmt"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

// Finding represents a compliance check result.
type Finding struct {
	Pass    bool
	Section string
	Detail  string
}

// CheckResult holds all findings from a compliance check.
type CheckResult struct {
	Findings []Finding
}

func (r *CheckResult) Add(pass bool, section, detail string) {
	r.Findings = append(r.Findings, Finding{Pass: pass, Section: section, Detail: detail})
}

func (r *CheckResult) PassCount() int {
	n := 0
	for _, f := range r.Findings {
		if f.Pass {
			n++
		}
	}
	return n
}

func (r *CheckResult) FailCount() int {
	return len(r.Findings) - r.PassCount()
}

// Check runs all compliance checks against the source directory.
func Check(specPath, sourceDir string) (*CheckResult, error) {
	s, err := spec.LoadSpec(specPath)
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	result := &CheckResult{}
	checkAPICoverage(s, sourceDir, result)
	checkModelConsistency(s, sourceDir, result)
	return result, nil
}
