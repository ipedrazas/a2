package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ipedrazas/a2/pkg/server"
	"github.com/spf13/cobra"
)

var (
	serverHost       string
	serverPort       int
	workspaceDir     string
	cleanupAfter     bool
	maxConcurrent    int
	cleanupInterval  time.Duration
	jobHistoryMaxAge time.Duration
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run a2 as a web server",
	Long: `Start a web server that provides an HTTP API and web UI for running checks.

The server accepts GitHub URLs, clones the repositories, runs checks, and returns results.
Ideal for CI/CD integration, code review automation, or on-demand analysis.`,
	RunE: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Server configuration
	serverCmd.Flags().StringVar(&serverHost, "host", "0.0.0.0", "Host to bind to")
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to listen on")
	serverCmd.Flags().StringVar(&workspaceDir, "workspace-dir", "/workspace/a2-cache", "Directory for cloned repositories")
	serverCmd.Flags().BoolVar(&cleanupAfter, "cleanup-after", true, "Clean up workspace after job completes")
	serverCmd.Flags().IntVar(&maxConcurrent, "max-concurrent", 5, "Maximum number of concurrent jobs")
	serverCmd.Flags().DurationVar(&cleanupInterval, "cleanup-interval", 1*time.Hour, "Interval for cleanup of old jobs/workspaces")
	serverCmd.Flags().DurationVar(&jobHistoryMaxAge, "job-history-max-age", 24*time.Hour, "Maximum age to keep job history")
}

func runServer(cmd *cobra.Command, args []string) error {
	// Create workspace manager
	wm := server.NewWorkspaceManager(workspaceDir, cleanupAfter)

	// Create job store
	jobStore := server.NewJobStore()

	// Create job queue
	queue := server.NewJobQueue(jobStore, maxConcurrent)

	// Create job processor (will be implemented in handlers.go)
	processor := func(ctx context.Context, job *server.Job) error {
		return server.ProcessJob(ctx, job, wm)
	}

	// Start the queue
	if err := queue.Start(processor); err != nil {
		return fmt.Errorf("failed to start job queue: %w", err)
	}
	defer queue.Stop()

	// Create HTTP server
	srv := server.NewServer(serverHost, serverPort, jobStore, queue, workspaceDir, cleanupAfter)

	// Start background cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startBackgroundCleanup(ctx, wm, jobStore, cleanupInterval, jobHistoryMaxAge)

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%s:%d", serverHost, serverPort)
		log.Printf("Starting a2 server on %s", addr)
		log.Printf("Workspace directory: %s", workspaceDir)
		log.Printf("Max concurrent jobs: %d", maxConcurrent)
		log.Printf("Cleanup after job: %v", cleanupAfter)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
		close(serverErr)
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("Received shutdown signal, gracefully shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		return nil
	case err := <-serverErr:
		return err
	}
}

func startBackgroundCleanup(ctx context.Context, wm *server.WorkspaceManager, js *server.JobStore, interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Clean up old jobs from store
			jobCount := js.CleanupOldJobs(maxAge)
			if jobCount > 0 {
				log.Printf("Cleaned up %d old jobs", jobCount)
			}

			// Clean up old workspace directories
			workspaceCount, err := wm.CleanupOld(maxAge)
			if err != nil {
				log.Printf("Error cleaning up workspaces: %v", err)
			} else if workspaceCount > 0 {
				log.Printf("Cleaned up %d old workspace directories", workspaceCount)
			}
		}
	}
}
