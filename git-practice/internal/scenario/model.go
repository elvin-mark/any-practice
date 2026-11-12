package scenario

// Scenario defines an interactive git drill or lab.
type Scenario struct {
	ID          string   `yaml:"id" json:"id"`
	Title       string   `yaml:"title" json:"title"`
	Category    string   `yaml:"category" json:"category"`
	Difficulty  string   `yaml:"difficulty" json:"difficulty"`
	Description string   `yaml:"description" json:"description"`
	Goal        string   `yaml:"goal" json:"goal"`
	Hints       []string `yaml:"hints" json:"hints"`
	Setup       string   `yaml:"setup" json:"setup"`
	Verify      string   `yaml:"verify" json:"verify"`
}
