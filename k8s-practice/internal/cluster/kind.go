package cluster

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/elvin-mark/any-practice/k8s-practice/internal/ui"
)

const (
	DefaultCluster1 = "cka-cluster-1"
	DefaultCluster2 = "cka-cluster-2"

	CalicoManifestURL = "https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml"
)

// KindConfig returns the Kind cluster YAML configuration
func KindConfig(clusterName string, multiNode bool, calico bool) string {
	var sb strings.Builder
	sb.WriteString("kind: Cluster\n")
	sb.WriteString("apiVersion: kind.x-k8s.io/v1alpha4\n")

	if calico {
		sb.WriteString("networking:\n")
		sb.WriteString("  disableDefaultCNI: true\n")
		sb.WriteString("  podSubnet: \"192.168.0.0/16\"\n")
	}

	sb.WriteString("nodes:\n")
	sb.WriteString("- role: control-plane\n")
	if multiNode {
		sb.WriteString("- role: worker\n")
	}

	return sb.String()
}

// ClusterStatus represents the state of a Kind cluster
type ClusterStatus struct {
	Name      string
	Context   string
	Running   bool
	NodeCount int
	Nodes     []NodeInfo
}

type NodeInfo struct {
	Name   string
	Role   string
	Status string
}

// CheckKindInstalled verifies kind and docker are available
func CheckPrerequisites() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is required but not found in PATH")
	}
	if _, err := exec.LookPath("kind"); err != nil {
		return fmt.Errorf("kind is required but not found in PATH")
	}
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl is required but not found in PATH")
	}
	return nil
}

// ListClusters returns all existing kind clusters
func ListClusters(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "kind", "get", "clusters")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing kind clusters: %w", err)
	}

	var clusters []string
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.Contains(l, "No kind clusters found") {
			clusters = append(clusters, l)
		}
	}
	return clusters, nil
}

// CreateCluster creates a kind cluster with specified configuration
func CreateCluster(ctx context.Context, name string, multiNode bool, calico bool) error {
	ui.PrintInfo("Creating Kind cluster '%s' (multiNode=%v, calico=%v)...", name, multiNode, calico)

	configFile, err := os.CreateTemp("", fmt.Sprintf("kind-config-%s-*.yaml", name))
	if err != nil {
		return fmt.Errorf("creating temp kind config: %w", err)
	}
	defer os.Remove(configFile.Name())

	configContent := KindConfig(name, multiNode, calico)
	if _, err := configFile.WriteString(configContent); err != nil {
		return fmt.Errorf("writing kind config: %w", err)
	}
	configFile.Close()

	cmd := exec.CommandContext(ctx, "kind", "create", "cluster", "--name", name, "--config", configFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running kind create cluster: %w", err)
	}

	if calico {
		ui.PrintInfo("Installing Calico CNI on cluster '%s'...", name)
		if err := InstallCalico(ctx, fmt.Sprintf("kind-%s", name)); err != nil {
			ui.PrintWarning("Failed to install Calico automatically: %v (you can apply manually via kubectl)", err)
		} else {
			ui.PrintSuccess("Calico CNI deployed to cluster '%s'", name)
		}
	}

	return nil
}

// InstallCalico deploys Calico CNI onto the specified cluster context
func InstallCalico(ctx context.Context, contextName string) error {
	cmd := exec.CommandContext(ctx, "kubectl", "--context", contextName, "apply", "-f", CalicoManifestURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("applying calico manifest: %s (%w)", stderr.String(), err)
	}

	ui.PrintInfo("Waiting for Calico node pods to initialize...")
	// Wait shortly for calico-node daemonset
	waitCmd := exec.CommandContext(ctx, "kubectl", "--context", contextName, "-n", "kube-system",
		"wait", "--for=condition=ready", "pod", "-l", "k8s-app=calico-node", "--timeout=90s")
	_ = waitCmd.Run()
	return nil
}

// DeleteCluster destroys a Kind cluster
func DeleteCluster(ctx context.Context, name string) error {
	ui.PrintInfo("Deleting Kind cluster '%s'...", name)
	cmd := exec.CommandContext(ctx, "kind", "delete", "cluster", "--name", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SwitchContext sets the current kube context
func SwitchContext(ctx context.Context, contextName string) error {
	if contextName == "" {
		return nil
	}

	// Check if context exists
	checkCmd := exec.CommandContext(ctx, "kubectl", "config", "get-contexts", contextName, "--no-headers")
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("context '%s' not found in kubeconfig", contextName)
	}

	cmd := exec.CommandContext(ctx, "kubectl", "config", "use-context", contextName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("switching context to '%s': %s (%w)", contextName, stderr.String(), err)
	}
	ui.PrintSuccess("Switched active Kubernetes context to: %s", ui.CyanBold(contextName))
	return nil
}

// GetCurrentContext returns the active kubectl context
func GetCurrentContext(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("retrieving current context: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetClusterDetails retrieves details about nodes in a cluster
func GetClusterDetails(ctx context.Context, contextName string) (*ClusterStatus, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "--context", contextName, "get", "nodes", "--no-headers")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getting nodes for context %s: %w", contextName, err)
	}

	status := &ClusterStatus{
		Context: contextName,
		Running: true,
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) >= 3 {
			status.Nodes = append(status.Nodes, NodeInfo{
				Name:   parts[0],
				Status: parts[1],
				Role:   parts[2],
			})
		}
	}
	status.NodeCount = len(status.Nodes)
	return status, nil
}
