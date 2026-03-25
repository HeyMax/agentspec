package diff

import (
	"fmt"
	"strings"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

// Impact describes the downstream effects of a Change.
type Impact struct {
	Change        Change
	AffectedAPIs  []string
	AffectedPages []string
	AffectedTests []string
	Suggestions   []string
}

// AnalyzeImpact traces each change through the spec's reference graph to find downstream effects.
func AnalyzeImpact(changes []Change, s *spec.FeatureSpec) []Impact {
	var impacts []Impact

	// Build reference indexes
	apiByBehavior := buildAPIByBehaviorIndex(s)
	apiByModel := buildAPIByModelIndex(s)
	pageByAPI := buildPageByAPIIndex(s)
	testByBehavior := buildTestByBehaviorIndex(s)
	testByAPI := buildTestByAPIIndex(s)

	for _, ch := range changes {
		impact := Impact{Change: ch}

		switch ch.Section {
		case "models":
			modelName := extractFirstSegment(ch.Path, "models.")
			// Model changed → find APIs that reference this model
			if apis, ok := apiByModel[modelName]; ok {
				impact.AffectedAPIs = apis
				// Trace to pages
				for _, apiID := range apis {
					if pages, ok := pageByAPI[apiID]; ok {
						impact.AffectedPages = appendUnique(impact.AffectedPages, pages...)
					}
				}
			}
			// Suggestions
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Update structs/interfaces for model %q", modelName),
				"Run migrations if field types or constraints changed",
			)

		case "behaviors":
			behID := extractFirstSegment(ch.Path, "behaviors.")
			// Behavior changed → find APIs linked to this behavior
			if apis, ok := apiByBehavior[behID]; ok {
				impact.AffectedAPIs = apis
				for _, apiID := range apis {
					if pages, ok := pageByAPI[apiID]; ok {
						impact.AffectedPages = appendUnique(impact.AffectedPages, pages...)
					}
				}
			}
			if tests, ok := testByBehavior[behID]; ok {
				impact.AffectedTests = tests
			}
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Review business logic for behavior %q", behID),
			)

		case "apis":
			apiID := extractFirstSegment(ch.Path, "apis.")
			impact.AffectedAPIs = []string{apiID}
			if pages, ok := pageByAPI[apiID]; ok {
				impact.AffectedPages = pages
			}
			if tests, ok := testByAPI[apiID]; ok {
				impact.AffectedTests = tests
			}
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Update API handler and client calls for %q", apiID),
			)

		case "ui":
			pageName := extractFirstSegment(ch.Path, "ui.pages.")
			impact.AffectedPages = []string{pageName}
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Update UI components for page %q", pageName),
			)

		case "acceptance":
			testID := extractFirstSegment(ch.Path, "acceptance.")
			impact.AffectedTests = []string{testID}
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Update test cases for %q", testID),
			)

		case "enums":
			enumName := extractFirstSegment(ch.Path, "enums.")
			// Enums can affect any model field that references them
			for modelName, model := range s.Models {
				for _, field := range model.Fields {
					if field.Type == "enum:"+enumName {
						if apis, ok := apiByModel[modelName]; ok {
							impact.AffectedAPIs = appendUnique(impact.AffectedAPIs, apis...)
						}
					}
				}
			}
			impact.Suggestions = append(impact.Suggestions,
				fmt.Sprintf("Update enum definitions for %q in all languages", enumName),
			)
		}

		impacts = append(impacts, impact)
	}

	return impacts
}

// --- Index builders ---

// buildAPIByBehaviorIndex maps behavior ID → []API IDs that reference it.
func buildAPIByBehaviorIndex(s *spec.FeatureSpec) map[string][]string {
	idx := make(map[string][]string)
	for apiID, api := range s.APIs {
		if api.Behavior != "" {
			idx[api.Behavior] = append(idx[api.Behavior], apiID)
		}
	}
	return idx
}

// buildAPIByModelIndex maps model name → []API IDs whose request/response reference that model.
func buildAPIByModelIndex(s *spec.FeatureSpec) map[string][]string {
	idx := make(map[string][]string)
	for apiID, api := range s.APIs {
		models := extractModelsFromAPI(api)
		for _, m := range models {
			idx[m] = append(idx[m], apiID)
		}
	}
	return idx
}

// buildPageByAPIIndex maps API ID → []page names that reference it.
func buildPageByAPIIndex(s *spec.FeatureSpec) map[string][]string {
	idx := make(map[string][]string)
	for pageName, page := range s.UI.Pages {
		if page.DataSource != "" {
			idx[page.DataSource] = append(idx[page.DataSource], pageName)
		}
		for _, comp := range page.Components {
			for _, action := range comp.Actions {
				if action.API != "" {
					idx[action.API] = append(idx[action.API], pageName)
				}
			}
		}
	}
	return idx
}

// buildTestByBehaviorIndex maps behavior ID → []acceptance test IDs.
func buildTestByBehaviorIndex(s *spec.FeatureSpec) map[string][]string {
	idx := make(map[string][]string)
	for _, ac := range s.Acceptance {
		for _, beh := range ac.Behaviors {
			idx[beh] = append(idx[beh], ac.ID)
		}
	}
	return idx
}

// buildTestByAPIIndex maps API ID → []acceptance test IDs.
func buildTestByAPIIndex(s *spec.FeatureSpec) map[string][]string {
	idx := make(map[string][]string)
	for _, ac := range s.Acceptance {
		for _, apiID := range ac.APIs {
			idx[apiID] = append(idx[apiID], ac.ID)
		}
	}
	return idx
}

// --- Helpers ---

// extractModelsFromAPI finds model names referenced in an API's request/response bodies.
func extractModelsFromAPI(api spec.API) []string {
	var models []string
	// Check response body type (often a model reference like "ref:Task")
	if api.Response.Success.Body != nil {
		if t := api.Response.Success.Body.Type; strings.HasPrefix(t, "ref:") {
			models = append(models, strings.TrimPrefix(t, "ref:"))
		}
		if api.Response.Success.Body.Items != nil {
			if t := api.Response.Success.Body.Items.Type; strings.HasPrefix(t, "ref:") {
				models = append(models, strings.TrimPrefix(t, "ref:"))
			}
		}
	}
	// Check request body fields for ref types
	if api.Request.Body != nil {
		for _, f := range api.Request.Body.Fields {
			if strings.HasPrefix(f.Type, "ref:") {
				models = append(models, strings.TrimPrefix(f.Type, "ref:"))
			}
		}
	}
	return models
}

// extractFirstSegment extracts the first path segment after a prefix.
// e.g. extractFirstSegment("models.Task.fields.name", "models.") → "Task"
func extractFirstSegment(path, prefix string) string {
	s := strings.TrimPrefix(path, prefix)
	if idx := strings.Index(s, "."); idx >= 0 {
		return s[:idx]
	}
	return s
}

func appendUnique(slice []string, items ...string) []string {
	seen := make(map[string]bool)
	for _, s := range slice {
		seen[s] = true
	}
	for _, item := range items {
		if !seen[item] {
			slice = append(slice, item)
			seen[item] = true
		}
	}
	return slice
}

// FormatImpacts produces a human-readable impact report.
func FormatImpacts(impacts []Impact) string {
	var b strings.Builder

	for _, imp := range impacts {
		b.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(string(imp.Change.Type)), imp.Change.Path))

		if len(imp.AffectedAPIs) > 0 {
			b.WriteString(fmt.Sprintf("  APIs affected: %s\n", strings.Join(imp.AffectedAPIs, ", ")))
		}
		if len(imp.AffectedPages) > 0 {
			b.WriteString(fmt.Sprintf("  Pages affected: %s\n", strings.Join(imp.AffectedPages, ", ")))
		}
		if len(imp.AffectedTests) > 0 {
			b.WriteString(fmt.Sprintf("  Tests affected: %s\n", strings.Join(imp.AffectedTests, ", ")))
		}
		if len(imp.Suggestions) > 0 {
			for _, s := range imp.Suggestions {
				b.WriteString(fmt.Sprintf("  → %s\n", s))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
