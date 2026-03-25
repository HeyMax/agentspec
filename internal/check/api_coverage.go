package check

import (
	"fmt"
	"strings"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

func checkAPICoverage(s *spec.FeatureSpec, sourceDir string, result *CheckResult) {
	scan, err := ScanDir(sourceDir)
	if err != nil {
		result.Add(false, "api-coverage", fmt.Sprintf("failed to scan source directory: %s", err))
		return
	}

	for apiName, api := range s.APIs {
		found := false
		for _, route := range scan.Routes {
			if route.Method == api.Method && pathMatches(api.Path, route.Path) {
				found = true
				break
			}
		}
		if found {
			result.Add(true, "api-coverage", fmt.Sprintf("%s %s (%s) — found in source", api.Method, api.Path, apiName))
		} else {
			result.Add(false, "api-coverage", fmt.Sprintf("%s %s (%s) — not found in source", api.Method, api.Path, apiName))
		}
	}
}

func pathMatches(specPath, sourcePath string) bool {
	specParts := strings.Split(strings.Trim(specPath, "/"), "/")
	sourceParts := strings.Split(strings.Trim(sourcePath, "/"), "/")
	if len(specParts) != len(sourceParts) {
		return false
	}
	for i := range specParts {
		sp := specParts[i]
		src := sourceParts[i]
		if isParam(sp) || isParam(src) {
			continue
		}
		if sp != src {
			return false
		}
	}
	return true
}

func isParam(segment string) bool {
	return (strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")) ||
		strings.HasPrefix(segment, ":")
}
