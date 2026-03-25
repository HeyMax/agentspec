package schema

import (
	"fmt"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

// Severity indicates how serious a diagnostic is.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarn
)

// Diagnostic represents a single validation finding.
type Diagnostic struct {
	Severity Severity
	Path     string
	Message  string
	Hint     string
}

func (d Diagnostic) IsError() bool {
	return d.Severity == SeverityError
}

// ValidationResult holds the results of a full validation run.
type ValidationResult struct {
	Diagnostics []Diagnostic
}

func (r *ValidationResult) AddError(path, msg, hint string) {
	r.Diagnostics = append(r.Diagnostics, Diagnostic{
		Severity: SeverityError,
		Path:     path,
		Message:  msg,
		Hint:     hint,
	})
}

func (r *ValidationResult) AddWarn(path, msg, hint string) {
	r.Diagnostics = append(r.Diagnostics, Diagnostic{
		Severity: SeverityWarn,
		Path:     path,
		Message:  msg,
		Hint:     hint,
	})
}

func (r *ValidationResult) ErrorCount() int {
	n := 0
	for _, d := range r.Diagnostics {
		if d.IsError() {
			n++
		}
	}
	return n
}

func (r *ValidationResult) WarnCount() int {
	return len(r.Diagnostics) - r.ErrorCount()
}

func (r *ValidationResult) HasErrors() bool {
	return r.ErrorCount() > 0
}

// Validate runs all validation passes on a loaded FeatureSpec.
func Validate(s *spec.FeatureSpec) *ValidationResult {
	result := &ValidationResult{}

	// Pass 1: structural validation
	validateStructure(s, result)

	// Pass 2: semantic rules
	validateSemanticRules(s, result)

	// Pass 3: cross-reference checks
	validateCrossReferences(s, result)

	return result
}

// validateStructure checks required top-level fields.
func validateStructure(s *spec.FeatureSpec, r *ValidationResult) {
	if s.Name == "" {
		r.AddError("spec.name", "spec name is required", "Add a 'name' field to the top of your spec")
	}

	if s.Kind != "" && s.Kind != "FeatureSpec" {
		r.AddError("spec.kind", fmt.Sprintf("unsupported kind %q, expected 'FeatureSpec'", s.Kind), "")
	}

	// Validate models
	for modelName, model := range s.Models {
		path := fmt.Sprintf("models.%s", modelName)
		if len(model.Fields) == 0 {
			r.AddError(path, "model has no fields defined", "Add at least one field to the model")
		}
		for fieldName, field := range model.Fields {
			fp := fmt.Sprintf("%s.fields.%s", path, fieldName)
			if field.Type == "" {
				r.AddError(fp, "field type is required", "Specify a type (string, int, bool, uuid, datetime, ref:X, enum:X)")
			}
		}
	}

	// Validate enums
	for enumName, enum := range s.Enums {
		path := fmt.Sprintf("enums.%s", enumName)
		if len(enum.Values) == 0 {
			r.AddError(path, "enum has no values defined", "Add at least one value to the enum")
		}
		for i, v := range enum.Values {
			if v.Name == "" {
				r.AddError(fmt.Sprintf("%s.values[%d]", path, i), "enum value name is required", "")
			}
		}
	}

	// Validate APIs
	for apiName, api := range s.APIs {
		path := fmt.Sprintf("apis.%s", apiName)
		if api.Method == "" {
			r.AddError(path+".method", "API method is required", "Use GET, POST, PUT, PATCH, or DELETE")
		}
		if api.Path == "" {
			r.AddError(path+".path", "API path is required", "Specify a URL path like /api/resource")
		}
	}
}
