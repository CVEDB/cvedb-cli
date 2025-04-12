package scripts

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

	FilePath   string
	ScriptID   string
	ScriptName string
}

var updateCfg = &UpdateConfig{}

func init() {
	ScriptsCmd.AddCommand(scriptsUpdateCmd)

	scriptsUpdateCmd.Flags().StringVar(&updateCfg.FilePath, "file", "", "YAML file for script definition")
	scriptsUpdateCmd.MarkFlagRequired("file")
	scriptsUpdateCmd.Flags().StringVar(&updateCfg.ScriptID, "id", "", "ID of the script to update")
	scriptsUpdateCmd.Flags().StringVar(&updateCfg.ScriptName, "name", "", "Name of the script to update")
}

var scriptsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a private script",
	Long:  `Update a private script by specifying either its ID or name. If neither is provided, the script name will be read from the YAML file.`,
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

	var scriptImportRequest cvedb.ScriptImport
	err = yaml.Unmarshal(data, &scriptImportRequest)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", cfg.FilePath, err)
	}
	var scriptID uuid.UUID
	if cfg.ScriptID != "" {
		scriptID, err = uuid.Parse(cfg.ScriptID)
		if err != nil {
			return fmt.Errorf("failed to parse script ID: %w", err)
		}
	} else {
		scriptName := cfg.ScriptName
		if scriptName == "" {
			scriptName = scriptImportRequest.Name
		}
		script, err := client.GetPrivateScriptByName(ctx, scriptName)
		if err != nil {
			return fmt.Errorf("failed to find script: %w", err)
		}
		scriptID = *script.ID
	}

	_, err = client.UpdatePrivateScript(ctx, &scriptImportRequest, scriptID)
	if err != nil {
		return fmt.Errorf("failed to update script: %w", err)
	}
	return nil
}
