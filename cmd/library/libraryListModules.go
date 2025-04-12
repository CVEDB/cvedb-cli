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

// libraryListModulesCmd represents the libraryListModules command
var libraryListModulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "List modules from the Cvedb library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Token = util.GetToken()
		cfg.BaseURL = util.Cfg.BaseUrl
		if err := runListModules(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	libraryListCmd.AddCommand(libraryListModulesCmd)
	libraryListModulesCmd.Flags().BoolVar(&cfg.JSONOutput, "json", false, "Display output in JSON format")
}

func runListModules(cfg *Config) error {
	client, err := cvedb.NewClient(
		cvedb.WithToken(cfg.Token),
		cvedb.WithBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	modules, err := client.ListLibraryModules(ctx)
	if err != nil {
		return fmt.Errorf("failed to get modules: %w", err)
	}

	if len(modules) == 0 {
		return fmt.Errorf("couldn't find any module in the library")
	}

	if cfg.JSONOutput {
		data, err := json.Marshal(modules)
		if err != nil {
			return fmt.Errorf("failed to marshal modules: %w", err)
		}
		fmt.Println(string(data))
	} else {
		err = display.PrintModules(os.Stdout, modules)
		if err != nil {
			return fmt.Errorf("failed to print modules: %w", err)
		}
	}
	return nil
}
