package diff

import (
	"fmt"
	"sort"

	"github.com/chenzhaohao/agentspec/internal/spec"
)

type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

type Change struct {
	Type    ChangeType
	Section string
	Path    string
	Detail  string
}

type DiffResult struct {
	Changes []Change
}

func (r *DiffResult) Add(typ ChangeType, section, path, detail string) {
	r.Changes = append(r.Changes, Change{Type: typ, Section: section, Path: path, Detail: detail})
}

func (r *DiffResult) HasChanges() bool {
	return len(r.Changes) > 0
}

func Diff(old, new *spec.FeatureSpec) *DiffResult {
	result := &DiffResult{}
	diffModels(old, new, result)
	diffEnums(old, new, result)
	diffBehaviors(old, new, result)
	diffAPIs(old, new, result)
	diffUI(old, new, result)
	return result
}

func DiffFiles(oldPath, newPath string) (*DiffResult, error) {
	oldSpec, err := spec.LoadSpec(oldPath)
	if err != nil {
		return nil, fmt.Errorf("loading old spec: %w", err)
	}
	newSpec, err := spec.LoadSpec(newPath)
	if err != nil {
		return nil, fmt.Errorf("loading new spec: %w", err)
	}
	return Diff(oldSpec, newSpec), nil
}

func diffModels(old, new *spec.FeatureSpec, r *DiffResult) {
	allKeys := mergeKeys(old.Models, new.Models)
	for _, name := range allKeys {
		oldModel, inOld := old.Models[name]
		newModel, inNew := new.Models[name]
		if !inOld {
			r.Add(Added, "models", "models."+name, fmt.Sprintf("model %s added", name))
			continue
		}
		if !inNew {
			r.Add(Removed, "models", "models."+name, fmt.Sprintf("model %s removed", name))
			continue
		}
		fieldKeys := mergeKeys(oldModel.Fields, newModel.Fields)
		for _, fn := range fieldKeys {
			oldField, fInOld := oldModel.Fields[fn]
			newField, fInNew := newModel.Fields[fn]
			path := fmt.Sprintf("models.%s.fields.%s", name, fn)
			if !fInOld {
				r.Add(Added, "models", path, fmt.Sprintf("field %s.%s added (type: %s)", name, fn, newField.Type))
			} else if !fInNew {
				r.Add(Removed, "models", path, fmt.Sprintf("field %s.%s removed", name, fn))
			} else if oldField.Type != newField.Type {
				r.Add(Modified, "models", path, fmt.Sprintf("field %s.%s type changed: %s → %s", name, fn, oldField.Type, newField.Type))
			} else if oldField.Required != newField.Required {
				r.Add(Modified, "models", path, fmt.Sprintf("field %s.%s required changed: %v → %v", name, fn, oldField.Required, newField.Required))
			}
		}
	}
}

func diffEnums(old, new *spec.FeatureSpec, r *DiffResult) {
	allKeys := mergeKeys(old.Enums, new.Enums)
	for _, name := range allKeys {
		_, inOld := old.Enums[name]
		newEnum, inNew := new.Enums[name]
		if !inOld {
			r.Add(Added, "enums", "enums."+name, fmt.Sprintf("enum %s added with %d values", name, len(newEnum.Values)))
		} else if !inNew {
			r.Add(Removed, "enums", "enums."+name, fmt.Sprintf("enum %s removed", name))
		}
	}
}

func diffBehaviors(old, new *spec.FeatureSpec, r *DiffResult) {
	allKeys := mergeKeys(old.Behaviors, new.Behaviors)
	for _, name := range allKeys {
		_, inOld := old.Behaviors[name]
		_, inNew := new.Behaviors[name]
		if !inOld {
			r.Add(Added, "behaviors", "behaviors."+name, fmt.Sprintf("behavior %s added", name))
		} else if !inNew {
			r.Add(Removed, "behaviors", "behaviors."+name, fmt.Sprintf("behavior %s removed", name))
		}
	}
}

func diffAPIs(old, new *spec.FeatureSpec, r *DiffResult) {
	allKeys := mergeKeys(old.APIs, new.APIs)
	for _, name := range allKeys {
		oldAPI, inOld := old.APIs[name]
		newAPI, inNew := new.APIs[name]
		path := "apis." + name
		if !inOld {
			r.Add(Added, "apis", path, fmt.Sprintf("API %s added (%s %s)", name, newAPI.Method, newAPI.Path))
		} else if !inNew {
			r.Add(Removed, "apis", path, fmt.Sprintf("API %s removed (%s %s)", name, oldAPI.Method, oldAPI.Path))
		} else if oldAPI.Method != newAPI.Method || oldAPI.Path != newAPI.Path {
			r.Add(Modified, "apis", path, fmt.Sprintf("API %s changed: %s %s → %s %s", name, oldAPI.Method, oldAPI.Path, newAPI.Method, newAPI.Path))
		}
	}
}

func diffUI(old, new *spec.FeatureSpec, r *DiffResult) {
	allKeys := mergeKeys(old.UI.Pages, new.UI.Pages)
	for _, name := range allKeys {
		_, inOld := old.UI.Pages[name]
		_, inNew := new.UI.Pages[name]
		path := "ui.pages." + name
		if !inOld {
			r.Add(Added, "ui", path, fmt.Sprintf("page %s added", name))
		} else if !inNew {
			r.Add(Removed, "ui", path, fmt.Sprintf("page %s removed", name))
		}
	}
}

func mergeKeys[V any](a, b map[string]V) []string {
	seen := map[string]bool{}
	for k := range a {
		seen[k] = true
	}
	for k := range b {
		seen[k] = true
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
