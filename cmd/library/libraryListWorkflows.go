package library

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/cvedb/cvedb-cli/pkg/display"
	"github.com/cvedb/cvedb-cli/util"
	"github.com/spf13/cobra"
)

// libraryListWorkflowsCmd represents the libraryListWorkflows command
var libraryListWorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List workflows from the Cvedb library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Token = util.GetToken()
		cfg.BaseURL = util.Cfg.BaseUrl
		if err := runListWorkflows(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	libraryListCmd.AddCommand(libraryListWorkflowsCmd)
	libraryListWorkflowsCmd.Flags().BoolVar(&cfg.JSONOutput, "json", false, "Display output in JSON format")
}

func runListWorkflows(cfg *Config) error {
	client, err := cvedb.NewClient(
		cvedb.WithToken(cfg.Token),
		cvedb.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	workflows, err := client.ListLibraryWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to get workflows: %w", err)
	}

	if len(workflows) == 0 {
		return fmt.Errorf("couldn't find any workflow in the library")
	}

	if cfg.JSONOutput {
		data, err := json.Marshal(workflows)
		if err != nil {
			return fmt.Errorf("failed to marshal workflows: %w", err)
		}
		fmt.Println(string(data))
	} else {
		err = display.PrintWorkflows(os.Stdout, workflows)
		if err != nil {
			return fmt.Errorf("failed to print workflows: %w", err)
		}
	}

	return nil
}
