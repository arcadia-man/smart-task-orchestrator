package docker

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "log"
    "os/exec"
    "strings"
    "time"
)

type Client struct {
    dockerHost string
}

type JobSpec struct {
    Image         string
    Command       string
    ContainerName string
    Env           map[string]string
    RunID         string
}

type ExecutionResult struct {
    ExitCode     int
    ContainerID  string
    ErrorMessage string
    Logs         []string
    LogStream    <-chan string
}

func NewClient(dockerHost string) (*Client, error) {
    // Test Docker connection
    cmd := exec.Command("docker", "version")
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("Docker is not available: %w", err)
    }

    log.Println("✅ Docker client initialized")
    return &Client{dockerHost: dockerHost}, nil
}

func (c *Client) ExecuteJob(ctx context.Context, spec *JobSpec) (*ExecutionResult, error) {
    log.Printf("🐳 Starting Docker execution: image=%s, container=%s", spec.Image, spec.ContainerName)

    // Create log stream channel
    logStream := make(chan string, 100)
    var logs []string

    // Pull image first
    if err := c.pullImage(ctx, spec.Image); err != nil {
        return nil, fmt.Errorf("failed to pull image: %w", err)
    }

    // Build docker run command
    args := []string{
        "run",
        "--rm", // Remove container after execution
        "--name", spec.ContainerName,
        "-i", // Interactive
    }

    // Add environment variables
    for key, value := range spec.Env {
        args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
    }

    // Add image and command
    args = append(args, spec.Image)

    // Split command into shell execution
    args = append(args, "/bin/sh", "-c", spec.Command)

    // Create command
    cmd := exec.CommandContext(ctx, "docker", args...)

    // Get stdout and stderr pipes
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
    }

    // Start the command
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("failed to start Docker container: %w", err)
    }

    // Get container ID
    containerID := spec.ContainerName

    // Stream logs from both stdout and stderr
    go c.streamLogs(stdout, "stdout", logStream, &logs)
    go c.streamLogs(stderr, "stderr", logStream, &logs)

    // Wait for command to complete
    err = cmd.Wait()

    // Close log stream
    close(logStream)

    // Get exit code
    exitCode := 0
    errorMessage := ""

    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            exitCode = exitError.ExitCode()
        } else {
            exitCode = -1
            errorMessage = err.Error()
        }
    }

    log.Printf("🐳 Docker execution completed: container=%s, exitCode=%d",
        spec.ContainerName, exitCode)

    return &ExecutionResult{
        ExitCode:     exitCode,
        ContainerID:  containerID,
        ErrorMessage: errorMessage,
        Logs:         logs,
        LogStream:    logStream,
    }, nil
}

func (c *Client) pullImage(ctx context.Context, image string) error {
    log.Printf("📥 Pulling Docker image: %s", image)

    cmd := exec.CommandContext(ctx, "docker", "pull", image)

    // Capture output for logging
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to pull image %s: %w\nOutput: %s", image, err, string(output))
    }

    log.Printf("✅ Image pulled successfully: %s", image)
    return nil
}

func (c *Client) streamLogs(reader io.Reader, source string, logStream chan<- string, logs *[]string) {
    scanner := bufio.NewScanner(reader)

    for scanner.Scan() {
        line := scanner.Text()

        // Add timestamp and source prefix
        timestampedLine := fmt.Sprintf("[%s] [%s] %s",
            time.Now().Format("15:04:05"), source, line)

        // Add to logs slice
        *logs = append(*logs, timestampedLine)

        // Send to stream (non-blocking)
        select {
        case logStream <- timestampedLine:
        default:
            // Channel full, skip this line
        }
    }

    if err := scanner.Err(); err != nil {
        errorLine := fmt.Sprintf("[%s] [%s] ERROR: %v",
            time.Now().Format("15:04:05"), source, err)

        *logs = append(*logs, errorLine)

        select {
        case logStream <- errorLine:
        default:
        }
    }
}

func (c *Client) StopContainer(ctx context.Context, containerName string) error {
    log.Printf("🛑 Stopping container: %s", containerName)

    cmd := exec.CommandContext(ctx, "docker", "stop", containerName)
    if err := cmd.Run(); err != nil {
        // Container might already be stopped
        log.Printf("⚠️  Failed to stop container %s: %v", containerName, err)
    }

    return nil
}

func (c *Client) RemoveContainer(ctx context.Context, containerName string) error {
    log.Printf("🗑️  Removing container: %s", containerName)

    cmd := exec.CommandContext(ctx, "docker", "rm", "-f", containerName)
    if err := cmd.Run(); err != nil {
        // Container might already be removed
        log.Printf("⚠️  Failed to remove container %s: %v", containerName, err)
    }

    return nil
}

func (c *Client) ListRunningContainers(ctx context.Context) ([]string, error) {
    cmd := exec.CommandContext(ctx, "docker", "ps", "--format", "{{.Names}}")

    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to list containers: %w", err)
    }

    lines := strings.Split(strings.TrimSpace(string(output)), "\n")
    var containers []string

    for _, line := range lines {
        if line != "" {
            containers = append(containers, line)
        }
    }

    return containers, nil
}
