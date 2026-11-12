package scenario

// Scenario defines an interactive SQL lab challenge.
type Scenario struct {
	ID          string   `yaml:"id" json:"id"`
	Title       string   `yaml:"title" json:"title"`
	Category    string   `yaml:"category" json:"category"`
	Difficulty  string   `yaml:"difficulty" json:"difficulty"`
	Description string   `yaml:"description" json:"description"`
	Goal        string   `yaml:"goal" json:"goal"`
	Hints       []string `yaml:"hints" json:"hints"`
	Setup       string   `yaml:"setup" json:"setup"`
	Solution    string   `yaml:"solution" json:"solution"`
	VerifyType  string   `yaml:"verify_type" json:"verify_type"` // "result_match" or "plan_match"
	VerifyPlan  string   `yaml:"verify_plan" json:"verify_plan"`
}
