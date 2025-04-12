package files

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

type ListConfig struct {
	Token   string
	BaseURL string

	SearchQuery string
	JSONOutput  bool
}

var listCfg = &ListConfig{}

func init() {
	FilesCmd.AddCommand(filesListCmd)

	filesListCmd.Flags().StringVar(&listCfg.SearchQuery, "query", "", "Filter listed files using the specified search query")
	filesListCmd.Flags().BoolVar(&listCfg.JSONOutput, "json", false, "Display output in JSON format")
}

// filesListCmd represents the filesGet command
var filesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in the Cvedb file storage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listCfg.Token = util.GetToken()
		listCfg.BaseURL = util.Cfg.BaseUrl
		if err := runList(listCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runList(cfg *ListConfig) error {
	client, err := cvedb.NewClient(
		cvedb.WithToken(cfg.Token),
		cvedb.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	files, err := client.SearchFiles(ctx, cfg.SearchQuery)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}

	if cfg.JSONOutput {
		data, err := json.Marshal(files)
		if err != nil {
			return fmt.Errorf("failed to marshall files: %w", err)
		}
		_, err = fmt.Fprintln(os.Stdout, string(data))
		if err != nil {
			return fmt.Errorf("failed to print files: %w", err)
		}
	} else {
		err = display.PrintFiles(os.Stdout, files)
		if err != nil {
			return fmt.Errorf("failed to print files: %w", err)
		}
	}
	return nil
}
