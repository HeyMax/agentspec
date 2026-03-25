package spec

import "time"

// FeatureSpec is the top-level structure of an agentspec YAML file.
type FeatureSpec struct {
	Kind        string                `yaml:"kind"`
	Version     string                `yaml:"version"`
	Name        string                `yaml:"name"`
	DisplayName string                `yaml:"display_name"`
	Owner       string                `yaml:"owner"`
	Description string                `yaml:"description"`
	Tags        []string              `yaml:"tags,omitempty"`
	Models      map[string]Model      `yaml:"models,omitempty"`
	Enums       map[string]Enum       `yaml:"enums,omitempty"`
	Behaviors   map[string]Behavior   `yaml:"behaviors,omitempty"`
	APIs        map[string]API        `yaml:"apis,omitempty"`
	UI          UISpec                `yaml:"ui,omitempty"`
	Acceptance  []AcceptanceCriteria  `yaml:"acceptance,omitempty"`

	// Internal metadata (not from YAML)
	FilePath  string    `yaml:"-"`
	LoadedAt  time.Time `yaml:"-"`
}

// ---------- Models ----------

type Model struct {
	Description string           `yaml:"description,omitempty"`
	Fields      map[string]Field `yaml:"fields"`
	States      *StateMachine    `yaml:"states,omitempty"`
	Indexes     []Index          `yaml:"indexes,omitempty"`
}

type Field struct {
	Type        string       `yaml:"type"`
	Required    bool         `yaml:"required,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Constraints *Constraints `yaml:"constraints,omitempty"`
	Default     interface{}  `yaml:"default,omitempty"`
}

type Constraints struct {
	MinLength *int     `yaml:"min_length,omitempty"`
	MaxLength *int     `yaml:"max_length,omitempty"`
	Min       *float64 `yaml:"min,omitempty"`
	Max       *float64 `yaml:"max,omitempty"`
	Pattern   string   `yaml:"pattern,omitempty"`
	Unique    bool     `yaml:"unique,omitempty"`
	Immutable bool     `yaml:"immutable,omitempty"`
}

type StateMachine struct {
	Field       string              `yaml:"field"`
	Values      []string            `yaml:"values"`
	Transitions map[string][]string `yaml:"transitions"`
}

type Index struct {
	Fields []string `yaml:"fields"`
	Unique bool     `yaml:"unique,omitempty"`
}

// ---------- Enums ----------

type Enum struct {
	Description string      `yaml:"description,omitempty"`
	Values      []EnumValue `yaml:"values"`
}

type EnumValue struct {
	Name        string      `yaml:"name"`
	Value       interface{} `yaml:"value"`
	Description string      `yaml:"description,omitempty"`
}

// ---------- Behaviors ----------

type Behavior struct {
	Description   string         `yaml:"description,omitempty"`
	Actor         string         `yaml:"actor,omitempty"`
	TargetModel   string         `yaml:"target_model,omitempty"`
	Preconditions []Condition    `yaml:"preconditions,omitempty"`
	Effects       []Effect       `yaml:"effects,omitempty"`
	Errors        []BehaviorError `yaml:"errors,omitempty"`
	CrossImpact   []CrossImpact  `yaml:"cross_impact,omitempty"`
	Idempotent    bool           `yaml:"idempotent,omitempty"`
	Async         bool           `yaml:"async,omitempty"`
}

type Condition struct {
	Description string `yaml:"description,omitempty"`
	Check       string `yaml:"check,omitempty"`
}

type Effect struct {
	Description string `yaml:"description,omitempty"`
	Mutation    string `yaml:"mutation,omitempty"`
}

type BehaviorError struct {
	Code       string `yaml:"code"`
	Condition  string `yaml:"condition,omitempty"`
	Message    string `yaml:"message"`
	HTTPStatus int    `yaml:"http_status,omitempty"`
}

type CrossImpact struct {
	Target      string `yaml:"target"`
	Effect      string `yaml:"effect,omitempty"`
	Description string `yaml:"description"`
}

// ---------- APIs ----------

type API struct {
	Method   string      `yaml:"method"`
	Path     string      `yaml:"path"`
	Summary  string      `yaml:"summary,omitempty"`
	Behavior string      `yaml:"behavior,omitempty"`
	Auth     *bool       `yaml:"auth,omitempty"`
	Request  APIRequest  `yaml:"request,omitempty"`
	Response APIResponse `yaml:"response,omitempty"`
}

func (a *API) RequiresAuth() bool {
	if a.Auth == nil {
		return true // default: requires auth
	}
	return *a.Auth
}

type APIRequest struct {
	Params map[string]ParamField `yaml:"params,omitempty"`
	Query  map[string]ParamField `yaml:"query,omitempty"`
	Body   *RequestBody          `yaml:"body,omitempty"`
}

type ParamField struct {
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type RequestBody struct {
	ContentType string                `yaml:"content_type,omitempty"`
	Fields      map[string]ParamField `yaml:"fields,omitempty"`
}

type APIResponse struct {
	Success SuccessResponse `yaml:"success,omitempty"`
	Errors  []ErrorResponse `yaml:"errors,omitempty"`
}

type SuccessResponse struct {
	Status     int                   `yaml:"status,omitempty"`
	Body       *ResponseBody         `yaml:"body,omitempty"`
}

type ResponseBody struct {
	Type       string                `yaml:"type,omitempty"`
	Fields     map[string]ParamField `yaml:"fields,omitempty"`
	Items      *ResponseItems        `yaml:"items,omitempty"`
	Pagination bool                  `yaml:"pagination,omitempty"`
}

type ResponseItems struct {
	Type string `yaml:"type"`
}

type ErrorResponse struct {
	Status      int    `yaml:"status"`
	Code        string `yaml:"code,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// ---------- UI ----------

type UISpec struct {
	Pages map[string]Page `yaml:"pages,omitempty"`
}

type Page struct {
	Type        string      `yaml:"type,omitempty"`
	Path        string      `yaml:"path"`
	Title       string      `yaml:"title,omitempty"`
	Description string      `yaml:"description,omitempty"`
	DataSource  string      `yaml:"data_source,omitempty"`
	Components  []Component `yaml:"components,omitempty"`
	Permissions *Permissions `yaml:"permissions,omitempty"`
}

type Component struct {
	ID      string                 `yaml:"id,omitempty"`
	Type    string                 `yaml:"type"`
	Props   map[string]interface{} `yaml:"props,omitempty"`
	Actions []Action               `yaml:"actions,omitempty"`
}

type Action struct {
	Trigger        string `yaml:"trigger"`
	API            string `yaml:"api,omitempty"`
	Confirm        bool   `yaml:"confirm,omitempty"`
	SuccessMessage string `yaml:"success_message,omitempty"`
	Redirect       string `yaml:"redirect,omitempty"`
}

type Permissions struct {
	Roles []string `yaml:"roles,omitempty"`
}

// ---------- Acceptance ----------

type AcceptanceCriteria struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Scenario  string   `yaml:"scenario,omitempty"`
	Behaviors []string `yaml:"behaviors,omitempty"`
	APIs      []string `yaml:"apis,omitempty"`
	Priority  string   `yaml:"priority,omitempty"`
}
