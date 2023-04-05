package store

import (
	"github.com/spf13/cobra"
)

// storeListCmd represents the storeList command
var storeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflows and tools from the CVEDB store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	StoreCmd.AddCommand(storeListCmd)
}
