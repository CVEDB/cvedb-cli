package files

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/cvedb/cvedb-cli/client/request"
	"github.com/cvedb/cvedb-cli/types"
	"github.com/cvedb/cvedb-cli/util"
)

var (
	Files string
)

// filesCmd represents the files command
var FilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage files in the Cvedb file storage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	FilesCmd.PersistentFlags().StringVar(&Files, "file", "", "File or files (comma-separated)")
	FilesCmd.MarkPersistentFlagRequired("file")

	FilesCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		_ = FilesCmd.Flags().MarkHidden("workflow")
		_ = FilesCmd.Flags().MarkHidden("project")
		_ = FilesCmd.Flags().MarkHidden("space")
		_ = FilesCmd.Flags().MarkHidden("url")

		command.Root().HelpFunc()(command, strings)
	})
}

func getMetadata(searchQuery string) ([]types.File, error) {
	resp := request.Cvedb.Get().DoF("file/?search=%s&vault=%s", searchQuery, util.GetVault())
	if resp == nil || resp.Status() != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.Status())
	}
	var metadata types.Files

	err := json.Unmarshal(resp.Body(), &metadata)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal file IDs response: %s", err)
	}

	return metadata.Results, nil
}
