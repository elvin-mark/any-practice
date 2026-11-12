package container

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerfilters "github.com/docker/docker/api/types/filters"
	dockerimage "github.com/docker/docker/api/types/image"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/elvin-mark/any-practice/linux-practice/internal/config"
	"github.com/elvin-mark/any-practice/linux-practice/internal/scenario"
)

const (
	LabelManagedBy = "linux-practice"
	LabelScenario  = "scenario-id"
)

// Manager wraps Docker API interactions for linux-practice.
type Manager struct {
	client *dockerclient.Client
}

// resolveDockerHost detects the appropriate Docker host socket across standard Docker,
// Docker Desktop for Linux, rootless Docker, and Podman.
func resolveDockerHost() string {
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		return host
	}

	home, err := os.UserHomeDir()
	if err == nil {
		desktopSock := home + "/.docker/desktop/docker.sock"
		if _, err := os.Stat(desktopSock); err == nil {
			return "unix://" + desktopSock
		}
	}

	uid := os.Getuid()
	userSock := fmt.Sprintf("/run/user/%d/docker.sock", uid)
	if _, err := os.Stat(userSock); err == nil {
		return "unix://" + userSock
	}

	podmanSock := fmt.Sprintf("/run/user/%d/podman/podman.sock", uid)
	if _, err := os.Stat(podmanSock); err == nil {
		return "unix://" + podmanSock
	}

	return ""
}

// NewManager creates and verifies a connection to the local Docker engine.
func NewManager() (*Manager, error) {
	opts := []dockerclient.Opt{
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	}

	if host := resolveDockerHost(); host != "" {
		opts = append(opts, dockerclient.WithHost(host))
	}

	cli, err := dockerclient.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("connecting to Docker daemon: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Docker daemon is unreachable. Make sure Docker/Podman is running: %w", err)
	}

	return &Manager{client: cli}, nil
}

// ContainerName returns standard container naming for a scenario.
func ContainerName(scenarioID string) string {
	return fmt.Sprintf("lp-%s", scenarioID)
}

// ImageExists checks whether an image tag is present in the local Docker image store.
func (m *Manager) ImageExists(ctx context.Context, imageName string) (bool, error) {
	images, err := m.client.ImageList(ctx, dockerimage.ListOptions{
		Filters: dockerfilters.NewArgs(dockerfilters.Arg("reference", imageName)),
	})
	if err != nil {
		return false, fmt.Errorf("checking image list: %w", err)
	}
	return len(images) > 0, nil
}

// BuildImage builds the base image from Dockerfile contents.
func (m *Manager) BuildImage(ctx context.Context, dockerfileContent []byte, imageName string) error {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0644,
		Size: int64(len(dockerfileContent)),
	}
	if err := tarWriter.WriteHeader(hdr); err != nil {
		return fmt.Errorf("writing tar header: %w", err)
	}
	if _, err := tarWriter.Write(dockerfileContent); err != nil {
		return fmt.Errorf("writing dockerfile to tar: %w", err)
	}
	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("closing tar writer: %w", err)
	}

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	resp, err := m.client.ImageBuild(ctx, &buf, buildOptions)
	if err != nil {
		return fmt.Errorf("triggering image build: %w", err)
	}
	defer resp.Body.Close()

	// Drain build output
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("reading image build output: %w", err)
	}

	return nil
}

// EnsureBaseImage makes sure the base image exists locally, building it if needed.
func (m *Manager) EnsureBaseImage(ctx context.Context, dockerfileContent []byte, imageName string) error {
	exists, err := m.ImageExists(ctx, imageName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	fmt.Printf("🔨 Building base systemd image %s (this only runs once)...\n", imageName)
	return m.BuildImage(ctx, dockerfileContent, imageName)
}

// StartScenario starts a container for the given scenario and runs setup.sh.
func (m *Manager) StartScenario(ctx context.Context, sc *scenario.Scenario) error {
	cName := ContainerName(sc.ID)

	// Clean up any stale container with the same name
	_ = m.StopScenario(ctx, sc.ID)

	portBindings := nat.PortMap{}
	for _, p := range sc.Container.Ports {
		parts := strings.Split(p, ":")
		if len(parts) == 2 {
			hostPort := parts[0]
			containerPort := parts[1]
			pKey := nat.Port(containerPort + "/tcp")
			portBindings[pKey] = []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: hostPort},
			}
		}
	}

	containerConfig := &dockercontainer.Config{
		Image: sc.Container.Image,
		Env:   sc.Container.Env,
		Labels: map[string]string{
			LabelManagedBy: "true",
			LabelScenario:  sc.ID,
		},
		StopSignal: "SIGRTMIN+3",
	}

	hostConfig := &dockercontainer.HostConfig{
		Privileged: sc.Container.Privileged,
		Tmpfs: map[string]string{
			"/run":      "rw,noexec,nosuid,size=65536k",
			"/run/lock": "rw,noexec,nosuid,size=65536k",
		},
		Binds: []string{
			"/sys/fs/cgroup:/sys/fs/cgroup:rw",
		},
		CgroupnsMode: "host",
		PortBindings: portBindings,
	}

	resp, err := m.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, cName)
	if err != nil {
		return fmt.Errorf("creating container %s: %w", cName, err)
	}

	if err := m.client.ContainerStart(ctx, resp.ID, dockercontainer.StartOptions{}); err != nil {
		_ = m.client.ContainerRemove(ctx, resp.ID, dockercontainer.RemoveOptions{Force: true})
		return fmt.Errorf("starting container %s: %w", cName, err)
	}

	// Wait for systemd to initialize
	if err := m.waitForSystemd(ctx, resp.ID); err != nil {
		_ = m.client.ContainerRemove(ctx, resp.ID, dockercontainer.RemoveOptions{Force: true})
		return fmt.Errorf("waiting for systemd: %w", err)
	}

	// Inject and run setup.sh inside the container
	if err := m.runScriptInside(ctx, resp.ID, "/tmp/setup.sh", sc.SetupScript); err != nil {
		_ = m.client.ContainerRemove(ctx, resp.ID, dockercontainer.RemoveOptions{Force: true})
		return fmt.Errorf("executing scenario setup.sh: %w", err)
	}

	return nil
}

// waitForSystemd polls until systemd has reached running or degraded state.
func (m *Manager) waitForSystemd(ctx context.Context, containerID string) error {
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		execCfg := types.ExecConfig{
			Cmd:          []string{"systemctl", "is-system-running"},
			AttachStdout: true,
			AttachStderr: true,
		}
		execID, err := m.client.ContainerExecCreate(ctx, containerID, execCfg)
		if err == nil {
			resp, err := m.client.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
			if err == nil {
				buf := new(strings.Builder)
				_, _ = io.Copy(buf, resp.Reader)
				resp.Close()

				status := strings.TrimSpace(buf.String())
				if strings.Contains(status, "running") ||
					strings.Contains(status, "degraded") ||
					strings.Contains(status, "maintenance") {
					return nil
				}
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(350 * time.Millisecond):
		}
	}
	return nil
}

// runScriptInside writes a script inside the container and executes it as root.
func (m *Manager) runScriptInside(ctx context.Context, containerID, destPath, content string) error {
	// Package script into tar archive
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: destPath,
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("tar header error: %w", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		return fmt.Errorf("tar write error: %w", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("tar close error: %w", err)
	}

	// Copy into container
	if err := m.client.CopyToContainer(ctx, containerID, "/", &buf, types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("copying script %s to container: %w", destPath, err)
	}

	// Execute script
	execCfg := types.ExecConfig{
		Cmd:          []string{"/bin/bash", destPath},
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := m.client.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return fmt.Errorf("creating exec process for %s: %w", destPath, err)
	}

	attachResp, err := m.client.ContainerExecAttach(ctx, execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return fmt.Errorf("attaching exec process for %s: %w", destPath, err)
	}
	defer attachResp.Close()

	var outBuf bytes.Buffer
	_, _ = stdcopy.StdCopy(&outBuf, &outBuf, attachResp.Reader)

	inspect, err := m.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return fmt.Errorf("inspecting exec process: %w", err)
	}

	if inspect.ExitCode != 0 {
		return fmt.Errorf("script failed with exit code %d:\n%s", inspect.ExitCode, outBuf.String())
	}

	return nil
}

// Verify executes verify.sh inside the container and returns success and output.
func (m *Manager) Verify(ctx context.Context, sc *scenario.Scenario) (bool, string, error) {
	cName := ContainerName(sc.ID)
	inspect, err := m.client.ContainerInspect(ctx, cName)
	if err != nil {
		return false, "", fmt.Errorf("container for %s is not running: %w", sc.ID, err)
	}
	if !inspect.State.Running {
		return false, "", fmt.Errorf("container for %s is stopped. Run 'lp start %s' first", sc.ID, sc.ID)
	}

	destPath := "/tmp/verify.sh"
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: destPath,
		Mode: 0755,
		Size: int64(len(sc.VerifyScript)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return false, "", fmt.Errorf("tar header error: %w", err)
	}
	if _, err := tw.Write([]byte(sc.VerifyScript)); err != nil {
		return false, "", fmt.Errorf("tar write error: %w", err)
	}
	if err := tw.Close(); err != nil {
		return false, "", fmt.Errorf("tar close error: %w", err)
	}

	if err := m.client.CopyToContainer(ctx, inspect.ID, "/", &buf, types.CopyToContainerOptions{}); err != nil {
		return false, "", fmt.Errorf("copying verify script to container: %w", err)
	}

	execCfg := types.ExecConfig{
		Cmd:          []string{"/bin/bash", destPath},
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := m.client.ContainerExecCreate(ctx, inspect.ID, execCfg)
	if err != nil {
		return false, "", fmt.Errorf("creating verify exec: %w", err)
	}

	attachResp, err := m.client.ContainerExecAttach(ctx, execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return false, "", fmt.Errorf("attaching verify exec: %w", err)
	}
	defer attachResp.Close()

	var outBuf bytes.Buffer
	_, _ = stdcopy.StdCopy(&outBuf, &outBuf, attachResp.Reader)

	execInspect, err := m.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return false, "", fmt.Errorf("inspecting verify process: %w", err)
	}

	passed := execInspect.ExitCode == 0
	cleanOutput := strings.TrimSpace(outBuf.String())
	return passed, cleanOutput, nil
}

// StopScenario stops and removes a scenario container.
func (m *Manager) StopScenario(ctx context.Context, scenarioID string) error {
	cName := ContainerName(scenarioID)
	_ = m.client.ContainerStop(ctx, cName, dockercontainer.StopOptions{})
	return m.client.ContainerRemove(ctx, cName, dockercontainer.RemoveOptions{Force: true})
}

// Attach drops the user into an interactive PTY shell inside the container as the student user.
func (m *Manager) Attach(scenarioID string) error {
	cName := ContainerName(scenarioID)
	cmd := exec.Command("docker", "exec", "-it", "-u", config.DefaultUser, "-w", "/home/"+config.DefaultUser, cName, "/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsRunning checks whether the scenario container is currently up.
func (m *Manager) IsRunning(ctx context.Context, scenarioID string) bool {
	cName := ContainerName(scenarioID)
	inspect, err := m.client.ContainerInspect(ctx, cName)
	if err != nil {
		return false
	}
	return inspect.State.Running
}
