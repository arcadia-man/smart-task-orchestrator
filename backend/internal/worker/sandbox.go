package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"smart-task-orchestrator/internal/models"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SandboxWorker struct {
	dockerCli *client.Client
	db        *mongo.Database
}

func NewSandboxWorker(db *mongo.Database) (*SandboxWorker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &SandboxWorker{dockerCli: cli, db: db}, nil
}

func (w *SandboxWorker) ExecuteJob(ctx context.Context, job models.Job, executionID string) error {
	execution := models.Execution{
		ID:        primitive.NewObjectID(),
		JobID:     job.ID,
		UserID:    job.UserID,
		Status:    models.StatusRunning,
		StartedAt: time.Now(),
	}

	w.db.Collection("executions").InsertOne(ctx, execution)

	// 1. Pull image (optimistic)
	// 2. Create container
	resp, err := w.dockerCli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: job.Image,
			Cmd:   []string{"sh", "-c", job.Command},
			Tty:   false,
		},
	})
	
	if err != nil {
		w.updateStatus(ctx, execution.ID, models.StatusFailed, err.Error())
		return fmt.Errorf("failed to create container: %w", err)
	}

	defer w.dockerCli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})

	// 3. Start container
	if _, err := w.dockerCli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		w.updateStatus(ctx, execution.ID, models.StatusFailed, err.Error())
		return fmt.Errorf("failed to start container: %w", err)
	}

	// 4. Capture logs
	out, err := w.dockerCli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err == nil {
		var logBuf bytes.Buffer
		io.Copy(&logBuf, out)
		out.Close()
		execution.Logs = logBuf.String()
	}

	// 5. Wait for container to exit
	waitResult := w.dockerCli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})
	select {
	case err := <-waitResult.Error:
		if err != nil {
			w.updateStatus(ctx, execution.ID, models.StatusFailed, err.Error())
			return err
		}
	case <-waitResult.Result:
	}

	w.updateStatus(ctx, execution.ID, models.StatusCompleted, execution.Logs)
	return nil
}

func (w *SandboxWorker) updateStatus(ctx context.Context, id primitive.ObjectID, status models.JobStatus, logs string) {
	now := time.Now()
	w.db.Collection("executions").UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status":   status,
			"logs":     logs,
			"ended_at": now,
		},
	})
}
