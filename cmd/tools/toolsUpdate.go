package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/cvedb/cvedb-cli/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type UpdateConfig struct {
	Token   string
	BaseURL string

	FilePath string
	ToolID   string
	ToolName string
}

var updateCfg = &UpdateConfig{}

func init() {
	ToolsCmd.AddCommand(toolsUpdateCmd)

	toolsUpdateCmd.Flags().StringVar(&updateCfg.FilePath, "file", "", "YAML file for tool definition")
	toolsUpdateCmd.MarkFlagRequired("file")
	toolsUpdateCmd.Flags().StringVar(&updateCfg.ToolID, "id", "", "ID of the tool to update")
	toolsUpdateCmd.Flags().StringVar(&updateCfg.ToolName, "name", "", "Name of the tool to update")
}

var toolsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a private tool integration",
	Long:  `Update a private tool integration by specifying either its ID or name. If neither is provided, the tool name will be read from the YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateCfg.Token = util.GetToken()
		updateCfg.BaseURL = util.Cfg.BaseUrl
		if err := runUpdate(updateCfg); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func runUpdate(cfg *UpdateConfig) error {
	data, err := os.ReadFile(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", cfg.FilePath, err)
	}

	client, err := cvedb.NewClient(cvedb.WithToken(cfg.Token), cvedb.WithBaseURL(cfg.BaseURL))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := context.Background()

	var toolImportRequest cvedb.ToolImport
	err = yaml.Unmarshal(data, &toolImportRequest)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", cfg.FilePath, err)
	}

	var toolID uuid.UUID
	if cfg.ToolID != "" {
		toolID, err = uuid.Parse(cfg.ToolID)
		if err != nil {
			return fmt.Errorf("failed to parse tool ID: %w", err)
		}
	} else {
		toolName := cfg.ToolName
		if toolName == "" {
			toolName = toolImportRequest.Name
		}
		tool, err := client.GetPrivateToolByName(ctx, toolName)
		if err != nil {
			return fmt.Errorf("failed to find tool: %w", err)
		}
		toolID = *tool.ID
	}

	_, err = client.UpdatePrivateTool(ctx, &toolImportRequest, toolID)
	if err != nil {
		return fmt.Errorf("failed to update tool: %w", err)
	}
	return nil
}
