package schema

import (
	"fmt"

	"github.com/HeyMax/agentspec/internal/spec"
)

// validateCrossReferences checks for orphaned definitions and inter-section consistency.
func validateCrossReferences(s *spec.FeatureSpec, r *ValidationResult) {
	checkOrphanedBehaviors(s, r)
	checkOrphanedModels(s, r)
	checkUIAPIReferences(s, r)
	checkAcceptanceReferences(s, r)
}

// checkOrphanedBehaviors warns about behaviors not referenced by any API.
func checkOrphanedBehaviors(s *spec.FeatureSpec, r *ValidationResult) {
	referenced := map[string]bool{}
	for _, api := range s.APIs {
		if api.Behavior != "" {
			referenced[api.Behavior] = true
		}
	}
	for _, ac := range s.Acceptance {
		for _, b := range ac.Behaviors {
			referenced[b] = true
		}
	}

	for name := range s.Behaviors {
		if !referenced[name] {
			r.AddWarn(fmt.Sprintf("behaviors.%s", name), "behavior is not referenced by any API or acceptance criteria", "Consider linking it to an API endpoint")
		}
	}
}

// checkOrphanedModels warns about models not referenced by any behavior, API response, or field.
func checkOrphanedModels(s *spec.FeatureSpec, r *ValidationResult) {
	referenced := map[string]bool{}

	// Collect model references from behaviors
	for _, beh := range s.Behaviors {
		if beh.TargetModel != "" {
			referenced[beh.TargetModel] = true
		}
		for _, ci := range beh.CrossImpact {
			referenced[ci.Target] = true
		}
	}

	// Collect model references from fields (ref:X)
	for _, model := range s.Models {
		for _, field := range model.Fields {
			collectModelRef(field.Type, referenced)
		}
	}

	// Collect from API responses
	for _, api := range s.APIs {
		if api.Response.Success.Body != nil {
			collectModelRef(api.Response.Success.Body.Type, referenced)
		}
	}

	for name := range s.Models {
		if !referenced[name] {
			r.AddWarn(fmt.Sprintf("models.%s", name), "model is not referenced by any behavior, field, or API response", "Ensure this model is used somewhere in your spec")
		}
	}
}

func collectModelRef(typ string, refs map[string]bool) {
	if len(typ) > 2 && typ[:2] == "[]" {
		typ = typ[2:]
	}
	if len(typ) > 4 && typ[:4] == "ref:" {
		refs[typ[4:]] = true
	}
}

// checkUIAPIReferences validates that UI page data_source references existing APIs.
func checkUIAPIReferences(s *spec.FeatureSpec, r *ValidationResult) {
	for pageName, page := range s.UI.Pages {
		path := fmt.Sprintf("ui.pages.%s", pageName)
		if page.DataSource != "" {
			if _, ok := s.APIs[page.DataSource]; !ok {
				r.AddWarn(path+".data_source", fmt.Sprintf("references undefined API %q", page.DataSource), "Ensure the API is defined in the apis section")
			}
		}

		for _, comp := range page.Components {
			for _, action := range comp.Actions {
				if action.API != "" {
					if _, ok := s.APIs[action.API]; !ok {
						r.AddWarn(fmt.Sprintf("%s.components.%s.actions", path, comp.ID), fmt.Sprintf("action references undefined API %q", action.API), "")
					}
				}
			}
		}
	}
}

// checkAcceptanceReferences validates that acceptance criteria reference existing behaviors and APIs.
func checkAcceptanceReferences(s *spec.FeatureSpec, r *ValidationResult) {
	for i, ac := range s.Acceptance {
		path := fmt.Sprintf("acceptance[%d](%s)", i, ac.ID)

		for _, bName := range ac.Behaviors {
			if _, ok := s.Behaviors[bName]; !ok {
				r.AddError(path+".behaviors", fmt.Sprintf("references undefined behavior %q", bName), "")
			}
		}

		for _, aName := range ac.APIs {
			if _, ok := s.APIs[aName]; !ok {
				r.AddError(path+".apis", fmt.Sprintf("references undefined API %q", aName), "")
			}
		}
	}
}
