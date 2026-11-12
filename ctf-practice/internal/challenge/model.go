package challenge

// Challenge represents a CTF problem/puzzle.
type Challenge struct {
	ID          string   `yaml:"id" json:"id"`
	Title       string   `yaml:"title" json:"title"`
	Category    string   `yaml:"category" json:"category"`
	Difficulty  string   `yaml:"difficulty" json:"difficulty"`
	Points      int      `yaml:"points" json:"points"`
	Description string   `yaml:"description" json:"description"`
	Hints       []string `yaml:"hints" json:"hints"`
	Flag        string   `yaml:"flag" json:"flag"`
	Files       []string `yaml:"files,omitempty" json:"files,omitempty"`
}
