package scenario

// Domain constants representing the official CKA exam domains
const (
	DomainClusterArch = "Cluster Architecture, Installation & Configuration"
	DomainWorkloads   = "Workloads & Scheduling"
	DomainServicesNet = "Services & Networking"
	DomainStorage     = "Storage"
	DomainTrouble     = "Troubleshooting"
)

// Metadata defines the metadata.yaml configuration for a scenario
type Metadata struct {
	ID              string   `yaml:"id"`
	Title           string   `yaml:"title"`
	Domain          string   `yaml:"domain"`
	Difficulty      string   `yaml:"difficulty"`
	Weight          int      `yaml:"weight"`
	TargetContext   string   `yaml:"targetContext"`
	TargetNamespace string   `yaml:"targetNamespace"`
	Hints           []string `yaml:"hints"`
}

// Scenario represents a complete CKA lab scenario
type Scenario struct {
	Metadata
	Path        string // Directory containing the scenario
	Description string // Markdown content from problem.md
	Solution    string // Markdown content from solution.md
}

// VerificationResult holds the output and status of verifying a scenario
type VerificationResult struct {
	Passed  bool
	Message string
	Output  string
	Checks  []CheckResult
}

// CheckResult represents an individual assertion result inside verification
type CheckResult struct {
	Name    string
	Passed  bool
	Details string
}
