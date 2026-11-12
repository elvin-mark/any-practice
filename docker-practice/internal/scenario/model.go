package scenario

// Official Docker Practice Domain Categories
const (
	DomainDockerfile      = "Dockerfile & Image Optimization"
	DomainLifecycle       = "Container Lifecycle & Resources"
	DomainStorage         = "Storage, Volumes & Persistence"
	DomainNetworking      = "Networking & Service Discovery"
	DomainCompose         = "Docker Compose Orchestration"
	DomainTroubleshooting = "Troubleshooting & Diagnostics"
	DomainSecurity        = "Container Security & Hardening"
)

// Metadata defines the metadata.yaml configuration for a scenario
type Metadata struct {
	ID         string   `yaml:"id"`
	Title      string   `yaml:"title"`
	Domain     string   `yaml:"domain"`
	Difficulty string   `yaml:"difficulty"`
	Weight     int      `yaml:"weight"`
	Tags       []string `yaml:"tags,omitempty"`
	Hints      []string `yaml:"hints"`
}

// Scenario represents a complete Docker practice scenario
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
