package context

// RoleConfig defines what sections of a spec a given role should see.
type RoleConfig struct {
	Name                 string
	Description          string
	IncludeModels        bool // full model details
	ModelsTypeOnly       bool // only types, no constraints
	IncludeEnums         bool
	IncludeBehaviors     bool
	BehaviorsSummaryOnly bool // only description, no pre/post
	IncludeCrossImpact   bool
	IncludeAPIs          bool
	APIsSummaryOnly      bool
	IncludeUI            bool
	IncludeAcceptance    bool
}

// Predefined role configurations.
var roles = map[string]RoleConfig{
	"backend": {
		Name:             "Backend",
		Description:      "Backend developer — models, behaviors, APIs, and cross-impact analysis",
		IncludeModels:    true,
		ModelsTypeOnly:   false,
		IncludeEnums:     true,
		IncludeBehaviors: true,
		IncludeCrossImpact: true,
		IncludeAPIs:      true,
		APIsSummaryOnly:  false,
		IncludeUI:        false,
		IncludeAcceptance: false,
	},
	"frontend": {
		Name:                 "Frontend",
		Description:          "Frontend developer — model types, API contracts, and UI specifications",
		IncludeModels:        true,
		ModelsTypeOnly:       true,
		IncludeEnums:         true,
		IncludeBehaviors:     true,
		BehaviorsSummaryOnly: true,
		IncludeCrossImpact:   false,
		IncludeAPIs:          true,
		APIsSummaryOnly:      false,
		IncludeUI:            true,
		IncludeAcceptance:    false,
	},
	"qa": {
		Name:             "QA",
		Description:      "Quality assurance — full specification visibility",
		IncludeModels:    true,
		ModelsTypeOnly:   false,
		IncludeEnums:     true,
		IncludeBehaviors: true,
		IncludeCrossImpact: true,
		IncludeAPIs:      true,
		APIsSummaryOnly:  false,
		IncludeUI:        true,
		IncludeAcceptance: true,
	},
	"pm": {
		Name:             "PM",
		Description:      "Product manager — behaviors, UI, and acceptance criteria",
		IncludeModels:    true,
		ModelsTypeOnly:   true,
		IncludeEnums:     false,
		IncludeBehaviors: true,
		IncludeCrossImpact: false,
		IncludeAPIs:      true,
		APIsSummaryOnly:  true,
		IncludeUI:        true,
		IncludeAcceptance: true,
	},
}

// GetRole returns the RoleConfig for the given role name.
// Returns the config and true if found, zero value and false otherwise.
func GetRole(name string) (RoleConfig, bool) {
	r, ok := roles[name]
	return r, ok
}

// RoleNames returns all valid role names.
func RoleNames() []string {
	return []string{"backend", "frontend", "qa", "pm"}
}
