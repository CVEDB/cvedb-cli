package get

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cvedb/cvedb-cli/pkg/config"
	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	display "github.com/cvedb/cvedb-cli/pkg/display/run"
	"github.com/cvedb/cvedb-cli/pkg/stats"
	"github.com/cvedb/cvedb-cli/util"
	"github.com/google/uuid"

	"github.com/spf13/cobra"
)

// Config holds the configuration for the get command
type Config struct {
	Token   string
	BaseURL string

	Watch                 bool
	IncludePrimitiveNodes bool
	IncludeTaskGroupStats bool
	JSONOutput            bool

	RunID   string
	RunSpec config.WorkflowRunSpec
}

var cfg = &Config{}

func init() {
	GetCmd.Flags().BoolVar(&cfg.Watch, "watch", false, "Watch the workflow execution if it's still running")
	GetCmd.Flags().BoolVar(&cfg.IncludePrimitiveNodes, "show-params", false, "Show parameters in the workflow tree")
	GetCmd.Flags().BoolVar(&cfg.IncludeTaskGroupStats, "analyze-task-groups", false, "Show detailed statistics for task groups, including task counts, status distribution, and duration analysis (min/max/median/outliers) (experimental)")
	GetCmd.Flags().StringVar(&cfg.RunID, "run", "", "Get the status of a specific run")
	GetCmd.Flags().BoolVar(&cfg.JSONOutput, "json", false, "Display output in JSON format")
}

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Displays status of a workflow",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Token = util.GetToken()
		cfg.BaseURL = util.Cfg.BaseUrl
		cfg.RunSpec = config.WorkflowRunSpec{
			RunID:        cfg.RunID,
			SpaceName:    util.SpaceName,
			ProjectName:  util.ProjectName,
			WorkflowName: util.WorkflowName,
			URL:          util.URL,
		}
		if err := run(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func run(cfg *Config) error {
	client, err := cvedb.NewClient(
		cvedb.WithToken(cfg.Token),
		cvedb.WithBaseURL(cfg.BaseURL),
	)

	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	runs, err := cfg.RunSpec.GetRuns(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}
	if len(runs) != 1 {
		return fmt.Errorf("expected 1 run, got %d", len(runs))
	}
	run := runs[0]

	err = displayRunDetails(ctx, client, &run, cfg)
	if err != nil {
		return fmt.Errorf("failed to handle run output: %w", err)
	}
	return nil
}

func displayRunDetails(ctx context.Context, client *cvedb.Client, run *cvedb.Run, cfg *Config) error {
	// Fetch the complete run details if the fleet information is missing
	// This happens when the run is retrieved from the workflow runs list which returns a simplified run object
	var err error
	if run.Fleet == nil {
		run, err = client.GetRun(ctx, *run.ID)
		if err != nil {
			return fmt.Errorf("failed to get run: %w", err)
		}
	}

	insights, err := client.GetRunSubJobInsights(ctx, *run.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Couldn't get the run insights: %s", err)
	} else {
		run.RunInsights = insights
	}

	averageDuration, err := client.GetWorkflowRunsAverageDuration(ctx, *run.WorkflowInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Couldn't calculate average duration: %s", err)
	} else {
		run.AverageDuration = &cvedb.Duration{Duration: averageDuration}
	}

	fleet, err := client.GetFleet(ctx, *run.Fleet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Couldn't get the fleet: %s", err)
	} else {
		run.FleetName = fleet.Name
	}

	version, err := client.GetWorkflowVersion(ctx, *run.WorkflowVersionInfo)
	if err != nil {
		return fmt.Errorf("failed to get workflow version: %w", err)
	}
	subjobs, err := client.GetSubJobs(ctx, *run.ID)
	if err != nil {
		return fmt.Errorf("failed to get subjobs: %w", err)
	}
	subjobs = cvedb.LabelSubJobs(subjobs, *version)

	ipAddresses, err := client.GetRunIPAddresses(ctx, *run.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Couldn't get the run IP addresses: %s", err)
	} else {
		run.IPAddresses = ipAddresses
	}

	if cfg.IncludeTaskGroupStats {
		for i := range subjobs {
			if subjobs[i].TaskGroup {
				childSubJobs, err := client.GetChildSubJobs(ctx, subjobs[i].ID)
				if err != nil {
					return fmt.Errorf("failed to get child subjobs: %w", err)
				}
				subjobs[i].Children = childSubJobs
			}
		}
	}

	if cfg.JSONOutput {
		var jsonRun *JSONRun
		if cfg.IncludeTaskGroupStats {
			taskGroupStatsMap := make(map[uuid.UUID]stats.TaskGroupStats)
			for _, subjob := range subjobs {
				if subjob.TaskGroup {
					taskGroupStatsMap[subjob.ID] = stats.CalculateTaskGroupStats(subjob)
				}
			}
			jsonRun = NewJSONRun(run, subjobs, taskGroupStatsMap)
		} else {
			jsonRun = NewJSONRun(run, subjobs, nil)
		}
		data, err := json.MarshalIndent(jsonRun, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal run data: %w", err)
		}
		output := string(data)
		fmt.Println(output)
	} else {
		if cfg.Watch {
			watcher, err := display.NewRunWatcher(
				client,
				*run.ID,
				display.WithWorkflowVersion(version),
				display.WithIncludePrimitiveNodes(cfg.IncludePrimitiveNodes),
				display.WithIncludeTaskGroupStats(cfg.IncludeTaskGroupStats),
			)
			if err != nil {
				return fmt.Errorf("failed to create run watcher: %w", err)
			}

			err = watcher.Watch(ctx)
			if err != nil {
				return fmt.Errorf("failed to watch run: %w", err)
			}
		} else {
			printer := display.NewRunPrinter(cfg.IncludePrimitiveNodes, os.Stdout)
			printer.PrintAll(run, subjobs, version, cfg.IncludeTaskGroupStats)
		}
	}
	return nil
}
