package util

import (
	"strings"
	"unicode"
)

// PascalCase converts kebab-case or snake_case to PascalCase.
// "create-task" → "CreateTask", "task_status" → "TaskStatus"
func PascalCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	return strings.Join(words, "")
}

// CamelCase converts kebab-case or snake_case to camelCase.
// "create-task" → "createTask"
func CamelCase(s string) string {
	p := PascalCase(s)
	if len(p) == 0 {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// KebabCase converts PascalCase or snake_case to kebab-case.
// "CreateTask" → "create-task"
func KebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '-')
		}
		result = append(result, unicode.ToLower(r))
	}
	return strings.ReplaceAll(string(result), "_", "-")
}

// SnakeCase converts kebab-case or PascalCase to snake_case.
func SnakeCase(s string) string {
	return strings.ReplaceAll(KebabCase(s), "-", "_")
}

func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	// Split on camelCase boundaries
	var words []string
	var current []rune
	for i, r := range s {
		if r == ' ' {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
			continue
		}
		if unicode.IsUpper(r) && i > 0 && len(current) > 0 {
			words = append(words, string(current))
			current = nil
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}
	return words
}

// GoType converts a spec type string to a Go type.
func GoType(specType string) string {
	if strings.HasPrefix(specType, "ref:") {
		return "*" + strings.TrimPrefix(specType, "ref:")
	}
	if strings.HasPrefix(specType, "enum:") {
		return PascalCase(strings.TrimPrefix(specType, "enum:"))
	}
	if strings.HasPrefix(specType, "[]") {
		inner := strings.TrimPrefix(specType, "[]")
		return "[]" + GoType(inner)
	}
	switch specType {
	case "string":
		return "string"
	case "int":
		return "int64"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	case "uuid":
		return "string"
	case "datetime":
		return "time.Time"
	default:
		return "interface{}"
	}
}

// TsType converts a spec type string to a TypeScript type.
func TsType(specType string) string {
	if strings.HasPrefix(specType, "ref:") {
		return strings.TrimPrefix(specType, "ref:")
	}
	if strings.HasPrefix(specType, "enum:") {
		return PascalCase(strings.TrimPrefix(specType, "enum:"))
	}
	if strings.HasPrefix(specType, "[]") {
		inner := strings.TrimPrefix(specType, "[]")
		return TsType(inner) + "[]"
	}
	switch specType {
	case "string", "uuid":
		return "string"
	case "int", "float":
		return "number"
	case "bool":
		return "boolean"
	case "datetime":
		return "string" // ISO 8601
	default:
		return "unknown"
	}
}
