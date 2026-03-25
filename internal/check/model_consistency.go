package check

import (
	"fmt"
	"strings"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

func checkModelConsistency(s *spec.FeatureSpec, sourceDir string, result *CheckResult) {
	scan, err := ScanDir(sourceDir)
	if err != nil {
		result.Add(false, "model-consistency", fmt.Sprintf("failed to scan source directory: %s", err))
		return
	}

	structMap := map[string]StructEntry{}
	for _, st := range scan.Structs {
		structMap[st.Name] = st
	}

	for modelName, model := range s.Models {
		st, found := structMap[modelName]
		if !found {
			result.Add(false, "model-consistency", fmt.Sprintf("model %s — struct not found in source", modelName))
			continue
		}

		result.Add(true, "model-consistency", fmt.Sprintf("model %s — struct found in %s", modelName, st.File))

		sourceFields := map[string]bool{}
		for _, f := range st.Fields {
			sourceFields[strings.ToLower(f)] = true
		}

		for fieldName := range model.Fields {
			normalizedField := strings.ToLower(fieldName)
			if sourceFields[normalizedField] {
				result.Add(true, "model-consistency", fmt.Sprintf("  %s.%s — found", modelName, fieldName))
			} else {
				result.Add(false, "model-consistency", fmt.Sprintf("  %s.%s — not found in struct", modelName, fieldName))
			}
		}
	}
}
