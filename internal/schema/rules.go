package schema

import (
	"fmt"
	"strings"

	"github.com/HeyMax/agentspec/internal/spec"
)

// validateSemanticRules checks type references, enum references, state machines, etc.
func validateSemanticRules(s *spec.FeatureSpec, r *ValidationResult) {
	checkTypeReferences(s, r)
	checkStateMachines(s, r)
	checkBehaviorReferences(s, r)
	checkAPIBehaviorLinks(s, r)
}

// checkTypeReferences validates ref: and enum: type prefixes point to existing definitions.
func checkTypeReferences(s *spec.FeatureSpec, r *ValidationResult) {
	for modelName, model := range s.Models {
		for fieldName, field := range model.Fields {
			path := fmt.Sprintf("models.%s.fields.%s", modelName, fieldName)
			checkFieldType(s, field.Type, path, r)
		}
	}

	// Check API request/response types
	for apiName, api := range s.APIs {
		basePath := fmt.Sprintf("apis.%s", apiName)

		if api.Request.Body != nil {
			for fieldName, f := range api.Request.Body.Fields {
				path := fmt.Sprintf("%s.request.body.%s", basePath, fieldName)
				checkFieldType(s, f.Type, path, r)
			}
		}

		if api.Response.Success.Body != nil && api.Response.Success.Body.Type != "" {
			checkFieldType(s, api.Response.Success.Body.Type, basePath+".response.success.body", r)
		}
	}
}

func checkFieldType(s *spec.FeatureSpec, typ, path string, r *ValidationResult) {
	// Handle array prefix
	if strings.HasPrefix(typ, "[]") {
		typ = strings.TrimPrefix(typ, "[]")
	}

	if strings.HasPrefix(typ, "ref:") {
		refName := strings.TrimPrefix(typ, "ref:")
		if _, ok := s.Models[refName]; !ok {
			r.AddError(path, fmt.Sprintf("references undefined model %q", refName), fmt.Sprintf("Define model %q in the models section or fix the reference", refName))
		}
	}

	if strings.HasPrefix(typ, "enum:") {
		enumName := strings.TrimPrefix(typ, "enum:")
		if _, ok := s.Enums[enumName]; !ok {
			r.AddError(path, fmt.Sprintf("references undefined enum %q", enumName), fmt.Sprintf("Define enum %q in the enums section or fix the reference", enumName))
		}
	}
}

// checkStateMachines validates state machine definitions.
func checkStateMachines(s *spec.FeatureSpec, r *ValidationResult) {
	for modelName, model := range s.Models {
		if model.States == nil {
			continue
		}
		sm := model.States
		path := fmt.Sprintf("models.%s.states", modelName)

		// State field must exist in model
		if sm.Field != "" {
			if _, ok := model.Fields[sm.Field]; !ok {
				r.AddError(path+".field", fmt.Sprintf("state field %q not found in model fields", sm.Field), "Add the field to the model or correct the state machine field name")
			}
		}

		// All transition sources must be declared values
		valueSet := map[string]bool{}
		for _, v := range sm.Values {
			valueSet[v] = true
		}

		for from, targets := range sm.Transitions {
			if !valueSet[from] {
				r.AddError(path+".transitions", fmt.Sprintf("transition source %q is not a declared state value", from), "Add it to the values list")
			}
			for _, to := range targets {
				if !valueSet[to] {
					r.AddError(path+".transitions", fmt.Sprintf("transition target %q is not a declared state value", to), "Add it to the values list")
				}
			}
		}

		// Warn about states with no transitions defined
		for _, v := range sm.Values {
			if _, ok := sm.Transitions[v]; !ok {
				r.AddWarn(path+".transitions", fmt.Sprintf("state %q has no transitions defined (terminal state?)", v), "If this is intentional, add an empty transition list")
			}
		}
	}
}

// checkBehaviorReferences validates behavior target_model references.
func checkBehaviorReferences(s *spec.FeatureSpec, r *ValidationResult) {
	for name, beh := range s.Behaviors {
		path := fmt.Sprintf("behaviors.%s", name)
		if beh.TargetModel != "" {
			if _, ok := s.Models[beh.TargetModel]; !ok {
				r.AddError(path+".target_model", fmt.Sprintf("references undefined model %q", beh.TargetModel), "")
			}
		}

		for i, ci := range beh.CrossImpact {
			if ci.Target != "" {
				if _, ok := s.Models[ci.Target]; !ok {
					r.AddWarn(fmt.Sprintf("%s.cross_impact[%d]", path, i), fmt.Sprintf("cross-impact target %q is not a defined model", ci.Target), "This may be intentional if the target is an external system")
				}
			}
		}
	}
}

// checkAPIBehaviorLinks ensures API.behavior references exist.
func checkAPIBehaviorLinks(s *spec.FeatureSpec, r *ValidationResult) {
	for name, api := range s.APIs {
		if api.Behavior != "" {
			if _, ok := s.Behaviors[api.Behavior]; !ok {
				r.AddError(fmt.Sprintf("apis.%s.behavior", name), fmt.Sprintf("references undefined behavior %q", api.Behavior), "Define the behavior or remove the reference")
			}
		}
	}
}
