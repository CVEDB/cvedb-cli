package display

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/google/uuid"
	"github.com/gosuri/uilive"
)

// RunWatcher handles watching and displaying the status of a workflow run
type RunWatcher struct {
	client                   *cvedb.Client
	runID                    uuid.UUID
	workflowVersion          *cvedb.WorkflowVersion
	includePrimitiveNodes    bool
	includeTaskGroupStats    bool
	ci                       bool
	writer                   *uilive.Writer
	mutex                    *sync.Mutex
	fetchedTaskGroupChildren map[uuid.UUID][]cvedb.SubJob // Cache to store completed children data
}

// RunWatcherOption is a function that configures a RunWatcher
type RunWatcherOption func(*RunWatcher)

// WithIncludePrimitiveNodes configures whether to include primitive nodes
func WithIncludePrimitiveNodes(include bool) RunWatcherOption {
	return func(w *RunWatcher) {
		w.includePrimitiveNodes = include
	}
}

// WithIncludeTaskGroupStats configures whether to include task group stats
func WithIncludeTaskGroupStats(include bool) RunWatcherOption {
	return func(w *RunWatcher) {
		w.includeTaskGroupStats = include
	}
}

// WithCI configures CI mode
func WithCI(ci bool) RunWatcherOption {
	return func(w *RunWatcher) {
		w.ci = ci
	}
}

// WithWorkflowVersion sets the workflow version for the watcher
func WithWorkflowVersion(version *cvedb.WorkflowVersion) RunWatcherOption {
	return func(w *RunWatcher) {
		w.workflowVersion = version
	}
}

// NewRunWatcher creates a new RunWatcher instance
func NewRunWatcher(client *cvedb.Client, runID uuid.UUID, opts ...RunWatcherOption) (*RunWatcher, error) {
	w := &RunWatcher{
		client:                   client,
		runID:                    runID,
		writer:                   uilive.New(),
		mutex:                    &sync.Mutex{},
		fetchedTaskGroupChildren: make(map[uuid.UUID][]cvedb.SubJob),
	}

	for _, opt := range opts {
		opt(w)
	}

	// If workflow version is not set, fetch it from the client
	if w.workflowVersion == nil {
		run, err := w.client.GetRun(context.Background(), w.runID)
		if err != nil {
			return nil, fmt.Errorf("failed to get run: %w", err)
		}
		if run == nil {
			return nil, fmt.Errorf("run not found")
		}
		if run.WorkflowVersionInfo == nil {
			return nil, fmt.Errorf("workflow version info not found in run")
		}
		version, err := w.client.GetWorkflowVersion(context.Background(), *run.WorkflowVersionInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow version: %w", err)
		}
		w.workflowVersion = version
	}

	return w, nil
}

// Watch starts watching the run and displaying its status
func (w *RunWatcher) Watch(ctx context.Context) error {
	w.writer.Start()
	defer w.writer.Stop()

	interruptErr := make(chan error, 1)
	go func() {
		interruptErr <- w.handleInterrupt(ctx)
	}()

	printer := NewRunPrinter(w.includePrimitiveNodes, w.writer)

	run, err := w.client.GetRun(ctx, w.runID)
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}

	fleet, err := w.client.GetFleet(ctx, *run.Fleet)
	if err != nil {
		return fmt.Errorf("failed to get fleet: %w", err)
	}
	fleetName := fleet.Name

	averageDuration, err := w.client.GetWorkflowRunsAverageDuration(ctx, *run.WorkflowInfo)
	if err != nil {
		return fmt.Errorf("failed to get average duration: %w", err)
	}

	for {
		select {
		case err := <-interruptErr:
			if err == nil || err.Error() == "execution interrupted by user" {
				return nil
			}
			return err
		default:
			w.mutex.Lock()
			run, err := w.client.GetRun(ctx, w.runID)
			if err != nil {
				w.mutex.Unlock()
				return fmt.Errorf("failed to get run: %w", err)
			}

			if run == nil {
				w.mutex.Unlock()
				return nil
			}

			subJobs, err := w.client.GetSubJobs(ctx, w.runID)
			if err != nil {
				w.mutex.Unlock()
				return fmt.Errorf("failed to get sub-jobs: %w", err)
			}

			if w.includeTaskGroupStats {
				for i := range subJobs {
					if subJobs[i].TaskGroup {
						// Only reload children if the task group is still running or hasn't been fetched before
						if subJobs[i].Status == "RUNNING" || len(w.fetchedTaskGroupChildren[subJobs[i].ID]) == 0 {
							childSubJobs, err := w.client.GetChildSubJobs(ctx, subJobs[i].ID)
							if err != nil {
								w.mutex.Unlock()
								return fmt.Errorf("failed to get child sub-jobs: %w", err)
							}
							w.fetchedTaskGroupChildren[subJobs[i].ID] = childSubJobs
						}
						subJobs[i].Children = w.fetchedTaskGroupChildren[subJobs[i].ID]
					}
				}
			}

			insights, err := w.client.GetRunSubJobInsights(ctx, w.runID)
			if err != nil {
				w.mutex.Unlock()
				return fmt.Errorf("failed to get run insights: %w", err)
			}
			run.RunInsights = insights
			run.FleetName = fleetName
			run.AverageDuration = &cvedb.Duration{Duration: averageDuration}

			printer.PrintAll(run, subJobs, w.workflowVersion, w.includeTaskGroupStats)
			_ = w.writer.Flush()

			if run.Finished {
				w.mutex.Unlock()
				return nil
			}

			w.mutex.Unlock()
			time.Sleep(time.Second)
		}
	}
}

// handleInterrupt handles the interrupt signal (Ctrl+C)
func (w *RunWatcher) handleInterrupt(ctx context.Context) error {
	defer w.mutex.Unlock()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	w.mutex.Lock()

	if w.ci {
		return w.client.StopRun(ctx, w.runID)
	} else {
		fmt.Println("The program will exit. Would you like to stop the remote execution? (Y/N)")
		var answer string
		for {
			_, _ = fmt.Scan(&answer)
			if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
				return w.client.StopRun(ctx, w.runID)
			} else if strings.ToLower(answer) == "n" || strings.ToLower(answer) == "no" {
				return fmt.Errorf("execution interrupted by user")
			}
		}
	}
}
