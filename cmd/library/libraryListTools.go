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

// libraryListToolsCmd represents the libraryListTools command
var libraryListToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List tools from the Cvedb library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Token = util.GetToken()
		cfg.BaseURL = util.Cfg.BaseUrl
		if err := runListTools(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	libraryListCmd.AddCommand(libraryListToolsCmd)
	libraryListToolsCmd.Flags().BoolVar(&cfg.JSONOutput, "json", false, "Display output in JSON format")
}

func runListTools(cfg *Config) error {
	client, err := cvedb.NewClient(
		cvedb.WithToken(cfg.Token),
		cvedb.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	tools, err := client.ListLibraryTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tools: %w", err)
	}

	if len(tools) == 0 {
		return fmt.Errorf("couldn't find any tool in the library")
	}

	if cfg.JSONOutput {
		data, err := json.Marshal(tools)
		if err != nil {
			return fmt.Errorf("failed to marshal tools: %w", err)
		}
		fmt.Println(string(data))
	} else {
		err = display.PrintTools(os.Stdout, tools)
		if err != nil {
			return fmt.Errorf("failed to print tools: %w", err)
		}
	}

	return nil
}
