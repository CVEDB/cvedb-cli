package get

import (
	"time"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/cvedb/cvedb-cli/pkg/stats"
	"github.com/google/uuid"
)

// JSONRun represents a simplified version of a Run for JSON output
type JSONRun struct {
	ID uuid.UUID `json:"id"`

	Status string `json:"status"`

	Author       string `json:"author"`
	CreationType string `json:"creation_type"`

	CreatedDate   *time.Time `json:"created_date"`
	StartedDate   *time.Time `json:"started_date"`
	CompletedDate *time.Time `json:"completed_date"`

	Duration        cvedb.Duration `json:"duration"`
	AverageDuration cvedb.Duration `json:"average_duration,omitempty"`

	WorkflowName        string    `json:"workflow_name"`
	WorkflowInfo        uuid.UUID `json:"workflow_info"`
	WorkflowVersionInfo uuid.UUID `json:"workflow_version_info"`

	Fleet        uuid.UUID `json:"fleet"`
	FleetName    string    `json:"fleet_name"`
	UseStaticIPs bool      `json:"use_static_ips"`
	Machines     int       `json:"machines"`
	IPAddresses  []string  `json:"ip_addresses"`

	RunInsights *cvedb.RunSubJobInsights `json:"run_insights,omitempty"`
	SubJobs     []JSONSubJob             `json:"subjobs"`
}

// JSONSubJob represents a simplified version of a SubJob for JSON output
type JSONSubJob struct {
	Label string `json:"label,omitempty"`
	Name  string `json:"name,omitempty"`

	Status  string `json:"status"`
	Message string `json:"message,omitempty"`

	StartedDate  *time.Time     `json:"started_date,omitempty"`
	FinishedDate *time.Time     `json:"finished_date,omitempty"`
	Duration     cvedb.Duration `json:"duration,omitempty"`

	IPAddress string `json:"ip_address,omitempty"`

	TaskGroup      bool                  `json:"task_group,omitempty"`
	TaskCount      int                   `json:"task_count,omitempty"`
	Children       []JSONSubJob          `json:"children,omitempty"`
	TaskIndex      int                   `json:"task_index,omitempty"`
	TaskGroupStats *stats.TaskGroupStats `json:"task_group_stats,omitempty"`
}

// NewJSONRun creates a new JSONRun from a cvedb.Run
func NewJSONRun(run *cvedb.Run, subjobs []cvedb.SubJob, taskGroupStatsMap map[uuid.UUID]stats.TaskGroupStats) *JSONRun {
	jsonRun := &JSONRun{
		ID:                  *run.ID,
		Status:              run.Status,
		Author:              run.Author,
		CreationType:        run.CreationType,
		CreatedDate:         run.CreatedDate,
		StartedDate:         run.StartedDate,
		CompletedDate:       run.CompletedDate,
		AverageDuration:     *run.AverageDuration,
		WorkflowName:        run.WorkflowName,
		WorkflowInfo:        *run.WorkflowInfo,
		WorkflowVersionInfo: *run.WorkflowVersionInfo,
		Fleet:               *run.Fleet,
		FleetName:           run.FleetName,
		UseStaticIPs:        *run.UseStaticIPs,
		IPAddresses:         run.IPAddresses,
		RunInsights:         run.RunInsights,
	}

	if run.Status == "RUNNING" {
		jsonRun.Duration = cvedb.Duration{Duration: time.Since(*run.StartedDate)}
	} else {
		jsonRun.Duration = cvedb.Duration{Duration: run.CompletedDate.Sub(*run.StartedDate)}
	}

	if run.Machines.Default != nil {
		jsonRun.Machines = *run.Machines.Default
	} else if run.Machines.SelfHosted != nil {
		jsonRun.Machines = *run.Machines.SelfHosted
	}

	jsonRun.SubJobs = make([]JSONSubJob, len(subjobs))
	for i, subjob := range subjobs {
		jsonRun.SubJobs[i] = *NewJSONSubJob(&subjob, taskGroupStatsMap)
	}

	return jsonRun
}

// NewJSONSubJob creates a new JSONSubJob from a cvedb.SubJob
func NewJSONSubJob(subjob *cvedb.SubJob, taskGroupStats map[uuid.UUID]stats.TaskGroupStats) *JSONSubJob {
	jsonSubJob := &JSONSubJob{
		Label:        subjob.Label,
		Name:         subjob.Name,
		Status:       subjob.Status,
		Message:      subjob.Message,
		StartedDate:  &subjob.StartedDate,
		FinishedDate: &subjob.FinishedDate,
		IPAddress:    subjob.IPAddress,
		TaskGroup:    subjob.TaskGroup,
		TaskIndex:    subjob.TaskIndex,
	}

	if !subjob.StartedDate.IsZero() && !subjob.FinishedDate.IsZero() {
		duration := subjob.FinishedDate.Sub(subjob.StartedDate)
		jsonSubJob.Duration = cvedb.Duration{Duration: duration}
	}

	if len(subjob.Children) > 0 {
		jsonSubJob.TaskCount = len(subjob.Children)
	}

	if subjob.TaskGroup && taskGroupStats != nil && stats.HasInterestingStats(subjob.Name) {
		stats, ok := taskGroupStats[subjob.ID]
		if ok {
			jsonSubJob.TaskGroupStats = &stats
		}
		for _, child := range subjob.Children {
			jsonSubJob.Children = append(jsonSubJob.Children, *NewJSONSubJob(&child, nil))
		}
	}

	return jsonSubJob
}
